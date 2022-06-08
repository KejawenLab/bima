package utils

import (
	"time"

	"github.com/gadelkareem/cachita"
)

type Cache struct {
	lifetime time.Duration
	pool     cachita.Cache
}

func NewCache(lifetime time.Duration) *Cache {
	return &Cache{
		lifetime: lifetime,
		pool:     cachita.Memory(),
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.pool.Put(key, &value, c.lifetime)
}

func (c *Cache) Invalidate(key string) {
	c.pool.Invalidate(key)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	var data interface{}
	err := c.pool.Get(key, &data)
	if err != nil {
		return nil, false
	}

	return data, true
}
