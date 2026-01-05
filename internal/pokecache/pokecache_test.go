package pokecache_test

import (
	"testing"
	"time"

	"pokedex/internal/pokecache"
)

func TestCache(t *testing.T) {
	cache := pokecache.NewCache(5 * time.Second)
    name := "Gopher"
	  _, found := cache.Get(name)
	  if found {
			t.Errorf("TestCache(Gopher) = true, want false")
		}
}
