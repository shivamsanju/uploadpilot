package plugins

import (
	"github.com/go-gorm/caches/v4"
)

func NewCachePlugin(cacher caches.Cacher) *caches.Caches {
	return &caches.Caches{Conf: &caches.Config{
		Cacher: cacher,
	}}
}
