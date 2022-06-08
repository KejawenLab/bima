package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Cache_Ttl(t *testing.T) {
	cache := NewCache(time.Millisecond * 100)
	cache.Set("test", "a")

	data, found := cache.Get("test")

	assert.True(t, found)
	assert.Equal(t, data.(string), "a")

	time.Sleep(200 * time.Millisecond)

	data, found = cache.Get("test")

	assert.False(t, found)
	assert.Nil(t, data)
}

func Test_Cache_Invalidate(t *testing.T) {
	cache := NewCache(time.Millisecond * 100)
	cache.Set("test", "a")

	data, found := cache.Get("test")

	assert.True(t, found)
	assert.Equal(t, data.(string), "a")

	cache.Invalidate("test")

	data, found = cache.Get("test")

	assert.False(t, found)
	assert.Nil(t, data)
}
