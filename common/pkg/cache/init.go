package cache

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/common/pkg/msg"
)

var (
	redisClient *redis.Client
)

func Init(redisAddr, redisPassword, redisUsername *string, redisTLS bool) error {
	opt := &redis.Options{
		Addr:         *redisAddr,
		Password:     *redisPassword,
		Username:     *redisUsername,
		DB:           0,
		PoolSize:     20, // Maximum number of connections in the pool
		MinIdleConns: 5,  // Minimum idle connections
	}

	if redisTLS {
		opt.TLSConfig = &tls.Config{}
	}

	rc := redis.NewClient(opt)

	if rc == nil || rc.Options() == nil {
		return fmt.Errorf(msg.RedisConnectionFailure, "failed to create redis client")
	}
	redisClient = rc

	res := rc.Ping(context.Background())
	if res == nil {
		return fmt.Errorf(msg.RedisConnectionFailure, "failed to ping redis")
	}

	p, err := res.Result()

	if err != nil || p == "" {
		return err
	}

	// invalidate all
	// if err := redisClient.FlushAll(context.Background()).Err(); err != nil {
	// 	return err
	// }

	infra.Log.Info(msg.RedisConnectionSuccess)
	return nil
}
