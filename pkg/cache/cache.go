// Package cache provides a high-throughput, sharded in-memory key-value cache.
package cache

import (
	"hash/fnv"
	"sync"
)

const defaultShardCount = 32

// Shard represents an isolated lock partition to prevent global mutex contention.
type Shard struct {
	mu    sync.RWMutex
	items map[string]any
}

// ShardedCache stores key-value pairs distributed across independent shards.
type ShardedCache struct {
	shardCount uint32
	shards     []*Shard
}

// NewShardedCache initializes a new ShardedCache with the default shard count.
func NewShardedCache() *ShardedCache {
	shards := make([]*Shard, defaultShardCount)
	for i := 0; i < defaultShardCount; i++ {
		shards[i] = &Shard{
			items: make(map[string]any),
		}
	}
	return &ShardedCache{
		shardCount: defaultShardCount,
		shards:     shards,
	}
}

func (c *ShardedCache) getShard(key string) *Shard {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return c.shards[h.Sum32()%c.shardCount]
}

// Get retrieves an item from the cache.
func (c *ShardedCache) Get(key string) (any, bool) {
	shard := c.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()

	val, found := shard.items[key]
	return val, found
}

// Set stores an item in the cache.
func (c *ShardedCache) Set(key string, value any) {
	shard := c.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.items[key] = value
}

// Delete removes an item from the cache.
func (c *ShardedCache) Delete(key string) {
	shard := c.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	delete(shard.items, key)
}

// Len returns the total number of items across all shards.
func (c *ShardedCache) Len() int {
	total := 0
	for _, shard := range c.shards {
		shard.mu.RLock()
		total += len(shard.items)
		shard.mu.RUnlock()
	}
	return total
}
