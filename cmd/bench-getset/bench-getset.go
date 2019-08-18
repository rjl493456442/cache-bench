package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/rjl493456442/cache-bench"
)

var (
	cacheName = flag.String("cache", "", "the name of testing cache(bigcache, fastcache, freecache)")

	keySize   = flag.Int("ksize", 32, "the size of entry key in bytes")
	valueSize = flag.Int("vsize", 100, "the size of entry value in bytes")
	cacheSize = flag.Int("size", 1024*1024*1024, "the size of cache in bytes")
	duration  = flag.Duration("duration", 5*time.Minute, "the duration of benchmark")
	mode      = flag.String("mode", "getset", "the mode of testing(get, set, getset")
	getp      = flag.Int("getp", 50, "the percentage of get operation if mode is getset")
)

func main() {
	flag.Parse()
	var (
		pureget  bool
		cacheTyp cache_bench.CacheTyp
		config   cache_bench.BenchConfig
	)
	switch strings.ToLower(*cacheName) {
	case "bigcache":
		cacheTyp = cache_bench.BigCache
	case "fastcache":
		cacheTyp = cache_bench.FastCache
	case "freecache":
		cacheTyp = cache_bench.FreeCache
	}
	if cacheTyp == cache_bench.Undefined {
		fatal("Undefined cache type")
	}
	var ops func(cnt uint64, cache cache_bench.Cache, key, value []byte) bool
	if *mode == "getset" {
		p := *getp
		ops = func(cnt uint64, cache cache_bench.Cache, key, value []byte) bool {
			if rand.Intn(100) < p {
				// get operation
				cache.Get(randomKey(key))
				return true
			} else {
				// set operation
				cache.Set(key, value)
				return false
			}
		}
	} else if *mode == "get" {
		pureget = true // initialization needed
		ops = func(cnt uint64, cache cache_bench.Cache, key, value []byte) bool {
			cache.Get(randomKey(key))
			return true
		}
	} else {
		ops = func(cnt uint64, cache cache_bench.Cache, key, value []byte) bool { cache.Set(key, value); return false }
	}
	config = cache_bench.BenchConfig{
		KeySize:   *keySize,
		ValueSize: *valueSize,
		CacheSize: *cacheSize,
		Duration:  *duration,
		Ops:       ops,
	}
	cache, err := cache_bench.NewCache(cacheTyp, *cacheSize, *valueSize)
	if err != nil {
		fatal("Failed to create cache", err)
	}
	// Initialize the cache if in pure get mode.
	var lastKey []byte
	if pureget {
		lastKey = cache_bench.Init(cache, config)
	}
	// Run benchmark and collect result as well as system metrics.
	fmt.Println(cache_bench.Bench(cache, config, lastKey))
}

func randomKey(key []byte) []byte {
	var ret = make([]byte, len(key))
	for i, k := range key {
		if k == 0 {
			ret[i] = byte(0)
		} else {
			ret[i] = byte(rand.Intn(int(k)))
		}
	}
	return ret
}

func fatal(msg ...interface{}) {
	fmt.Println(msg)
	os.Exit(1)
}
