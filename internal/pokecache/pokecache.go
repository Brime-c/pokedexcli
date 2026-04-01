package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	value     []byte
}

type Cache struct {
	entry    map[string]cacheEntry
	interval time.Duration
	mu       sync.Mutex
}

func NewCache(inter time.Duration) *Cache {
	cache := &Cache{
		entry:    make(map[string]cacheEntry),
		interval: inter,
	}
	go cache.reapLoop()
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entry[key] = cacheEntry{
		createdAt: time.Now(),
		value:     val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.entry[key]
	if !ok {
		return nil, false
	}
	return val.value, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.mu.Lock()
		for key, value := range c.entry {
			if time.Since(value.createdAt) > c.interval {
				delete(c.entry, key)
			}
		}
		c.mu.Unlock()
	}
}
