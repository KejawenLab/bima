package utils

import (
	"time"

	"github.com/allegro/bigcache/v3"
)

type Cache struct {
	pool *bigcache.BigCache
}

func NewCache(lifetime time.Duration) *Cache {
	cache, _ := bigcache.NewBigCache(bigcache.DefaultConfig(lifetime))

	return &Cache{pool: cache}
}

func (c *Cache) Set(key string, value []byte) {
	c.pool.Set(key, value)
}

func (c *Cache) Invalidate(key string) {
	c.pool.Delete(key)
}

func (c *Cache) Get(key string) ([]byte, bool) {
	data, err := c.pool.Get(key)
	if err != nil {
		return nil, false
	}

	return data, true
}
