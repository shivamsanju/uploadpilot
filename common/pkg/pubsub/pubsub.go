package pubsub

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/common/pkg/msg"
	commonutils "github.com/uploadpilot/uploadpilot/common/pkg/utils"
)

type RedisConfig struct {
	Addr     *string
	Password *string
	Username *string
	TLS      bool
}

type EventBus[T any] struct {
	event       string
	client      *redis.Client
	ctx         context.Context
	group       string
	consumerKey string
	wg          *sync.WaitGroup
}

func NewEventBus[T any](event, consumerKey string, c *RedisConfig) *EventBus[T] {
	opt := &redis.Options{
		Addr:     *c.Addr,
		Password: *c.Password,
		Username: *c.Username,
	}

	if c.TLS {
		opt.TLSConfig = &tls.Config{}
	}

	client := redis.NewClient(opt)

	return &EventBus[T]{
		event:       event,
		client:      client,
		ctx:         context.Background(),
		consumerKey: consumerKey,
		wg:          &sync.WaitGroup{},
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

func (bus *EventBus[T]) Subscribe(consumerGroup string, handler func(msg *T) error) {
	bus.group = consumerGroup
	err := bus.client.XGroupCreateMkStream(context.Background(), bus.event, consumerGroup, "$").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		infra.Log.Errorf("error creating redis stream: %s", err.Error())
	}
	infra.Log.Infof("created redis stream %s for group %s", bus.event, bus.group)

	messageChan := make(chan *redis.XMessage)

	go bus.listen(messageChan)
	go bus.startWorker(messageChan, handler)

}

func (bus *EventBus[T]) startWorker(messageChan chan *redis.XMessage, handler func(msg *T) error) {
	defer commonutils.Recover()
	for msg := range messageChan {
		bus.processMessage(handler, msg)
	}
}

func (bus *EventBus[T]) listen(messageChan chan *redis.XMessage) {
	defer commonutils.Recover()

	for {
		pendingMessages := bus.getPendingMessages()
		for _, message := range pendingMessages {
			messageChan <- &message
		}

		streams, err := bus.client.XReadGroup(bus.ctx, &redis.XReadGroupArgs{
			Group:    bus.group,
			Consumer: bus.consumerKey,
			Streams:  []string{bus.event, ">"},
			Count:    1,
			Block:    5 * time.Second,
		}).Result()

		if err == redis.Nil {
			continue
		} else if err != nil {
			infra.Log.Errorf(msg.ErrReadingFromStream, err.Error())
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
		infra.Log.Errorf(msg.ErrInvalidMessageFormat, message.ID)
		return
	}

	var m T
	if err := json.Unmarshal([]byte(data), &m); err != nil {
		infra.Log.Errorf(msg.ErrDecodingMessage, message.ID, err.Error())
		return
	}

	if err := handler(&m); err != nil {
		infra.Log.Errorf(msg.ErrorInHandler, message.ID, err.Error())
		return
	}

}

// Retry pending (unacknowledged) messages
func (bus *EventBus[T]) getPendingMessages() []redis.XMessage {
	pending, err := bus.client.XPending(bus.ctx, bus.event, bus.group).Result()

	if err != nil {
		infra.Log.Errorf(msg.ErrReadingFromStream, err.Error())
		return nil
	}

	if pending.Count <= 0 {
		return nil
	}

	infra.Log.Infof(msg.RetryingFailedMessages, pending.Count)
	messages, err := bus.client.XClaim(bus.ctx, &redis.XClaimArgs{
		Stream:   bus.event,
		Group:    bus.group,
		Consumer: bus.consumerKey,
		MinIdle:  10 * time.Second,
		Messages: []string{pending.Lower},
	}).Result()

	if err != nil {
		infra.Log.Errorf(msg.ErrClaimingMessage, err.Error())
		return nil
	}

	return messages
}

func (bus *EventBus[T]) Unsubscribe() {
	bus.client.Conn().Close()
}
