// Package cache provides an in-memory key-value cache.
package cache

import (
	"sync"
)

// GlobalMutexCache replaces sharding with a single global mutex.
type GlobalMutexCache struct {
	mu    sync.Mutex
	items map[string]any
}

// NewGlobalMutexCache creates a new cache instance using a global mutex.
func NewGlobalMutexCache() *GlobalMutexCache {
	return &GlobalMutexCache{
		items: make(map[string]any),
	}
}

// Get retrieves an item from the cache.
func (c *GlobalMutexCache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, found := c.items[key]
	return val, found
}

// Set stores an item in the cache.
func (c *GlobalMutexCache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = value
}

// Delete removes an item from the cache.
func (c *GlobalMutexCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}
