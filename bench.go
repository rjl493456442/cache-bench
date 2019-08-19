package cache_bench

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

type BenchConfig struct {
	KeySize   int                                                                 `json:"keysize"`   // The size of each key written
	ValueSize int                                                                 `json:"valuesize"` // The size of each value written
	CacheSize int                                                                 `json:"cachesize"` // The size of cache assigned
	Duration  time.Duration                                                       `json:"duration"`  // The duration of benchmark
	Ops       func(cnt uint64, cache Cache, key, value []byte, query []byte) bool `json:"-"`         // The specified testing operation
}

type Report struct {
	TPS       uint64 // Operations per second
	MemStat   runtime.MemStats
	CacheStat map[string]interface{}
}

// String returns the string representation of benchmark report.
func (r Report) String() string {
	var report string
	report += fmt.Sprintf("**************** Benchmark Report ****************\n")
	report += fmt.Sprintf("[1] TPS %d\n", r.TPS)
	report += fmt.Sprintf("[2] MemStat Alloc: %dMB, TotalAlloc:%dMB, Sys:%dMB, Mallocs:%d, Free:%d\n",
		toMegaBytes(r.MemStat.Alloc), toMegaBytes(r.MemStat.TotalAlloc), toMegaBytes(r.MemStat.Sys),
		r.MemStat.Mallocs, r.MemStat.Frees)
	report += fmt.Sprintf("[3] CacheStat %v\n", r.CacheStat)
	return report
}

func Init(cache Cache, config BenchConfig) []byte {
	var (
		key   = make([]byte, config.KeySize)
		value = make([]byte, config.ValueSize)
	)
	for i := 0; i < config.CacheSize/(config.KeySize+config.ValueSize); i++ {
		cache.Set(key, value)
		key, value = incBytes(key), incBytes(value)
	}
	return key
}

func Bench(cache Cache, config BenchConfig, startKey []byte) Report {
	var (
		cnt     uint64
		memstat runtime.MemStats
		query   = make([]byte, config.KeySize)
		key     = make([]byte, config.KeySize)
		value   = make([]byte, config.ValueSize)
		stop    = make(chan struct{})
	)
	if startKey != nil {
		copy(key, startKey)
	}
	time.AfterFunc(config.Duration, func() { close(stop) })

loop:
	for {
		select {
		case <-stop:
			break loop
		default:
		}
		if !config.Ops(cnt, cache, key, value, query) {
			// Increase key value pair if set operation is made
			key = incBytes(key)
			value = incBytes(value)
		}
		cnt += 1
		randomKey(query, key)
	}
	runtime.ReadMemStats(&memstat)
	return Report{
		TPS:       cnt / uint64(config.Duration.Seconds()),
		MemStat:   memstat,
		CacheStat: cache.Stat(),
	}
}

func incBytes(bytes []byte) []byte {
	for i := 0; i < len(bytes); i++ {
		bytes[i]++
		if bytes[i] == 0 {
			continue
		}
		return bytes
	}
	return bytes
}

func toMegaBytes(b uint64) uint64 {
	return b / 1024 / 1024
}

func randomKey(dst []byte, key []byte) {
	if len(dst) != len(key) {
		panic("invalid byte slice for random")
	}
	for i, k := range key {
		if k == 0 {
			dst[i] = byte(0)
		} else {
			dst[i] = byte(rand.Intn(int(k)))
		}
	}
}
