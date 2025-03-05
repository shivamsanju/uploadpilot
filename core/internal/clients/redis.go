package clients

import (
	"crypto/tls"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisOpts struct {
	Addr     string
	Username string
	Password string
	TLS      bool
}

func NewRedisClient(opts *RedisOpts) (*redis.Client, error) {
	if opts.Addr == "" {
		return nil, fmt.Errorf("redis address is required")
	}
	redisOpts := &redis.Options{
		Addr:     opts.Addr,
		Username: opts.Username,
		Password: opts.Password,
	}
	if opts.TLS {
		redisOpts.TLSConfig = &tls.Config{}
	}

	redisClient := redis.NewClient(redisOpts)
	return redisClient, nil
}
