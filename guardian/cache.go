package guardian

import (
	"time"

	"github.com/shaj13/libcache"
)

func setupCache() (*libcache.Cache, error) {
	cache := libcache.FIFO.New(0)
	cache.SetTTL(time.Minute * 5)
	cache.RegisterOnExpired(func(key, _ interface{}) {
		cache.Peek(key)
	})
	return &cache, nil
}
