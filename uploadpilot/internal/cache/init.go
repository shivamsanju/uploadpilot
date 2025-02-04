package cache

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/internal/infra"
)

var (
	redisClient *redis.Client
)

func Init() error {
	opt := &redis.Options{
		Addr:         config.RedisAddr,
		Password:     config.RedisPassword,
		Username:     config.RedisUsername,
		DB:           0,
		PoolSize:     20, // Maximum number of connections in the pool
		MinIdleConns: 5,  // Minimum idle connections
	}

	if config.RedisTLS {
		opt.TLSConfig = &tls.Config{}
	}

	rc := redis.NewClient(opt)

	if rc == nil || rc.Options() == nil {
		return fmt.Errorf("failed to create redis client")
	}
	redisClient = rc

	res := rc.Ping(context.Background())
	if res == nil {
		return fmt.Errorf("failed to ping redis")
	}

	p, err := res.Result()

	if err != nil || p == "" {
		return err
	}

	// invalidate all
	if err := redisClient.FlushAll(context.Background()).Err(); err != nil {
		return err
	}

	infra.Log.Info("successfully connected to redis!")
	return nil
}
