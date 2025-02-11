package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/uploadpilot/uploadpilot/common/pkg/db/dbutils"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/common/pkg/msg"
)

type Client[T any] struct {
	client *redis.Client
}

func NewClient[T any](defaultTTL time.Duration) *Client[T] {
	return &Client[T]{
		client: redisClient,
	}
}

func (r *Client[T]) Mutate(ctx context.Context, key string, dependentKeys []string, value T, dbMutateFn func(T) error, ttl time.Duration) error {
	err := dbMutateFn(value)
	if err != nil {
		return err
	}

	r.Set(ctx, key, value, ttl)
	Invalidate(ctx, dependentKeys...)

	return nil
}

func (r *Client[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf(msg.ErrFailedToMarshal, err.Error())
	}

	err = r.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		return fmt.Errorf(msg.RedisErrFailedToSet, err.Error())
	}

	return nil
}

func (r *Client[T]) Query(ctx context.Context, key string, result T, dbFetchFn func(T) error) error {
	data, err := r.client.Get(ctx, key).Bytes()

	if err != nil {
		if err == redis.Nil {
			// infra.Log.Infof("cache miss: %s", key)
		}
		err := dbFetchFn(result)
		if err != nil {
			return dbutils.DBError(err)
		}
		r.Set(ctx, key, result, 0)
		return nil
	}
	err = json.Unmarshal(data, result)

	return err
}

func Invalidate(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		err := redisClient.Del(ctx, key).Err()
		if err != nil {
			infra.Log.Errorf(msg.RedisErrFailedToDel, err.Error())
			return errors.New(msg.ErrUnknown)
		}
	}
	return nil
}
