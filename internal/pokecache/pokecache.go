package pokecache

import (
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
}

type cacheEntry struct {
	createdAt time.Time 
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	// fmt.Println("Initialize New Cache")
	c := Cache{
		cache: make(map[string]cacheEntry),
	}

	c.ReapLoop(interval)

	return &c
}

func (c *Cache) ReapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			for key, entry := range c.cache {
				age := time.Since(entry.createdAt)
				if age > interval { 
					// fmt.Println("found one that should be reaped")
					delete(c.cache, key)
				}
			}
		}
	}()
}

func (c *Cache) Get(url string) ([]byte, bool) {
	// fmt.Print("Get()")
	entry, ok := c.cache[url]
	if ok {
		// fmt.Println("Cache key found")
		return entry.val, true
	} 
	// fmt.Println("Cache miss")
	return nil, false
}

func (c *Cache) Add(url string, body []byte) {
	c.cache[url] = cacheEntry{
		createdAt: time.Now().UTC(),
		val: body,
	}

	// fmt.Println("Key added to cache")
}

func pokecache() {
	return	
}
