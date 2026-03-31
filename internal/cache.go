package pokecache

import (
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	value     []byte
}

type Cache struct {
	entry map[string]cacheEntry
}
