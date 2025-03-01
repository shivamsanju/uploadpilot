package cacheplugins

import (
	"github.com/go-gorm/caches/v4"
	"github.com/redis/go-redis/v9"
	"github.com/uploadpilot/manager/internal/db/cache"
)

func NewRedisCachesPlugin(redisClient *redis.Client) *caches.Caches {
	return &caches.Caches{Conf: &caches.Config{
		Cacher: cache.NewRedisCacher(redisClient),
	}}
}
