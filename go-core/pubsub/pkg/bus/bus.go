package bus

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func Recover() {
	if r := recover(); r != nil {
		zap.L().Sugar().Errorf("recovered from panic: %s", r)
	}
}

type EventBus[T any] struct {
	event  string
	client *redis.Client
	ctx    context.Context
	group  string
	key    string
	wg     *sync.WaitGroup
}

func NewEventBus[T any](event string, client *redis.Client) *EventBus[T] {
	return &EventBus[T]{
		event:  event,
		client: client,
		ctx:    context.Background(),
		wg:     &sync.WaitGroup{},
	}
}

func (bus *EventBus[T]) Publish(msg *T) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = bus.client.XAdd(bus.ctx, &redis.XAddArgs{
		Stream: bus.event,
		Values: map[string]interface{}{"data": data},
		ID:     "*",
	}).Result()

	return err
}

func (bus *EventBus[T]) Subscribe(consumerGroup string, consumerKey string, handler func(msg *T) error) error {
	bus.group = consumerGroup
	bus.key = consumerKey
	err := bus.client.XGroupCreateMkStream(context.Background(), bus.event, consumerGroup, "$").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return fmt.Errorf("error creating redis stream: %s", err.Error())
	}

	messageChan := make(chan *redis.XMessage)

	go bus.listen(messageChan)
	go bus.startWorker(messageChan, handler)

	return nil
}

func (bus *EventBus[T]) startWorker(messageChan chan *redis.XMessage, handler func(msg *T) error) {
	defer Recover()
	for msg := range messageChan {
		bus.processMessage(handler, msg)
	}
}

func (bus *EventBus[T]) listen(messageChan chan *redis.XMessage) {
	defer Recover()

	for {
		pendingMessages := bus.getPendingMessages()
		for _, message := range pendingMessages {
			messageChan <- &message
		}

		streams, err := bus.client.XReadGroup(bus.ctx, &redis.XReadGroupArgs{
			Group:    bus.group,
			Consumer: bus.key,
			Streams:  []string{bus.event, ">"},
			Count:    1,
			Block:    5 * time.Second,
		}).Result()

		if err == redis.Nil {
			continue
		} else if err != nil {
			zap.L().Sugar().Errorf("error reading from stream: %s", err.Error())
			continue
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				messageChan <- &message
			}
		}
	}
}

func (bus *EventBus[T]) processMessage(handler func(msg *T) error, message *redis.XMessage) {
	data, ok := message.Values["data"].(string)
	defer bus.client.XAck(bus.ctx, bus.event, bus.group, message.ID)
	if !ok {
		zap.L().Sugar().Errorf("message [%s] has an invalid message format", message.ID)
		return
	}

	var m T
	if err := json.Unmarshal([]byte(data), &m); err != nil {
		zap.L().Sugar().Errorf("error decoding message [%s]: %s", message.ID, err.Error())
		return
	}

	if err := handler(&m); err != nil {
		zap.L().Sugar().Errorf("error while handling message [%s]: %s", message.ID, err.Error())
		return
	}

}

// Retry pending (unacknowledged) messages
func (bus *EventBus[T]) getPendingMessages() []redis.XMessage {
	pending, err := bus.client.XPending(bus.ctx, bus.event, bus.group).Result()

	if err != nil {
		zap.L().Sugar().Errorf("error reading from stream: %s", err.Error())
		return nil
	}

	if pending.Count <= 0 {
		return nil
	}

	zap.L().Sugar().Infof("retrying failed messages: %d", pending.Count)
	messages, err := bus.client.XClaim(bus.ctx, &redis.XClaimArgs{
		Stream:   bus.event,
		Group:    bus.group,
		Consumer: bus.key,
		MinIdle:  10 * time.Second,
		Messages: []string{pending.Lower},
	}).Result()

	if err != nil {
		zap.L().Sugar().Errorf("error claiming message: %s", err.Error())
		return nil
	}

	return messages
}

func (bus *EventBus[T]) Unsubscribe() {
	bus.client.Conn().Close()
}
