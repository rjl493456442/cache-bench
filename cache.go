package cache_bench

import (
	"time"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/allegro/bigcache"
	"github.com/coocood/freecache"
)

// Cache is the interface of all GC friendly in-memory cache which contains
// necessary functions for testing purposes.
type Cache interface {
	// Set stores (k, v) in the cache.
	Set(k, v []byte)

	// Get retrieves the entry value specified by given key.
	Get(k []byte) []byte

	// Stat retrieves all cache statistic maintained by underlying cache
	// itself.
	Stat() map[string]interface{}
}

// CacheTyp is the type indicator of specific in-memory cache implementation.
type CacheTyp int

const (
	Undefined CacheTyp = iota
	BigCache
	FastCache
	FreeCache
)

// String returns the string representation of cache implementation.
func (t CacheTyp) String() string {
	switch t {
	case BigCache:
		return "BigCache"
	case FastCache:
		return "FastCache"
	case FreeCache:
		return "FreeCache"
	default:
		return "Undefined"
	}
}

// NewCache initializes a cache instance with given type and relative config.
func NewCache(t CacheTyp, size int, valuesize int) (Cache, error) {
	var (
		err   error
		cache Cache
	)
	switch t {
	case FastCache:
		cache, err = newFastCache(size, valuesize), nil
	case BigCache:
		cache, err = newBigCache(size)
	case FreeCache:
		cache, err = newFreeCache(size), nil
	}
	return cache, err
}

type bigCache struct {
	cache *bigcache.BigCache
}

func newBigCache(size int) (*bigCache, error) {
	config := bigcache.DefaultConfig(10 * time.Minute)
	config.HardMaxCacheSize = size / 1024 / 1024
	config.Verbose = false
	cache, err := bigcache.NewBigCache(config)
	if err != nil {
		return nil, err
	}
	return &bigCache{cache: cache}, nil
}

func (c *bigCache) Set(k, v []byte) {
	c.cache.Set(string(k), v)
}

func (c *bigCache) Get(k []byte) []byte {
	value, _ := c.cache.Get(string(k))
	return value
}

func (c *bigCache) Stat() map[string]interface{} {
	stats := make(map[string]interface{})
	stat := c.cache.Stats()
	stats["hits"] = stat.Hits
	stats["misses"] = stat.Misses
	stats["delHits"] = stat.DelHits
	stats["delMisses"] = stat.DelMisses
	stats["collisions"] = stat.Collisions
	return stats
}

type fastCache struct {
	valuebuf []byte
	cache    *fastcache.Cache
}

func newFastCache(size int, valuesize int) *fastCache {
	return &fastCache{
		cache:    fastcache.New(size),
		valuebuf: make([]byte, valuesize),
	}
}

func (c *fastCache) Set(k, v []byte) {
	c.cache.Set(k, v)
}

func (c *fastCache) Get(k []byte) []byte {
	return c.cache.Get(c.valuebuf[:0], k)
}

func (c *fastCache) Stat() map[string]interface{} {
	return nil
}

type freeCache struct {
	cache *freecache.Cache
}

func newFreeCache(size int) *freeCache { return &freeCache{freecache.NewCache(size)} }

func (c *freeCache) Set(k, v []byte) {
	c.cache.Set(k, v, 0) // no expire
}

func (c *freeCache) Get(k []byte) []byte {
	value, _ := c.cache.Get(k)
	return value
}

func (c *freeCache) Stat() map[string]interface{} {
	stats := make(map[string]interface{})
	stats["count"] = c.cache.EntryCount()
	stats["avgAccess"] = c.cache.AverageAccessTime()
	stats["hit"] = c.cache.HitCount()
	stats["miss"] = c.cache.MissCount()
	stats["hitRate"] = c.cache.HitRate()
	stats["overWrite"] = c.cache.OverwriteCount()
	return stats
}
