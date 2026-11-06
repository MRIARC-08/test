package cache

import (
	"fmt"
	"sync"
	"testing"
)

func TestShardedCache(t *testing.T) {
	c := NewShardedCache()

	c.Set("alpha", 123)
	val, ok := c.Get("alpha")
	if !ok || val != 123 {
		t.Fatalf("expected 123, got %v", val)
	}

	c.Delete("alpha")
	_, ok = c.Get("alpha")
	if ok {
		t.Fatalf("expected alpha to be deleted")
	}
}

func TestConcurrentAccess(t *testing.T) {
	c := NewShardedCache()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", id)
			c.Set(key, id)
			val, ok := c.Get(key)
			if !ok || val != id {
				t.Errorf("expected %d, got %v", id, val)
			}
		}(i)
	}
	wg.Wait()
}
