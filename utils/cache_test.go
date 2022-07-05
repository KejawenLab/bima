package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Cache_Ttl(t *testing.T) {
	cache := NewCache(time.Second * 1)
	cache.Set("test", []byte("a"))

	data, found := cache.Get("test")

	assert.True(t, found)
	assert.Equal(t, string(data), "a")

	time.Sleep(1 * time.Second)

	data, found = cache.Get("test")

	assert.False(t, found)
	assert.Nil(t, data)
}

func Test_Cache_Invalidate(t *testing.T) {
	cache := NewCache(time.Millisecond * 100)
	cache.Set("test", []byte("a"))

	data, found := cache.Get("test")

	assert.True(t, found)
	assert.Equal(t, string(data), "a")

	cache.Invalidate("test")

	data, found = cache.Get("test")

	assert.False(t, found)
	assert.Nil(t, data)
}
