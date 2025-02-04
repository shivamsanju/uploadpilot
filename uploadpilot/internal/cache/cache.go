package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/uploadpilot/uploadpilot/internal/infra"
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
	r.Invalidate(ctx, dependentKeys...)

	return nil
}

func (r *Client[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	err = r.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set value in cache: %w", err)
	}

	return nil
}

func (r *Client[T]) Query(ctx context.Context, key string, result T, dbFetchFn func(T) error) error {
	data, err := r.client.Get(ctx, key).Bytes()

	if err != nil {
		if err == redis.Nil {
			infra.Log.Infof("cache miss: %s", key)
		}
		infra.Log.Infof(("Before fetch ws: %+v"), result)
		err := dbFetchFn(result)
		if err != nil {
			return err
		}
		infra.Log.Infof(("After fetch ws: %+v"), result)
		r.Set(ctx, key, result, 0)
		return nil
	}
	err = json.Unmarshal(data, result)

	return err
}

func (r *Client[T]) Invalidate(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		err := r.client.Del(ctx, key).Err()
		if err != nil {
			infra.Log.Errorf("failed to invalidate cache: %s", err.Error())
			return errors.New("there was an issue processing your request. please try again later")
		}
	}
	return nil
}
