---
title: 'Cache bench'
disqus: hackmd
---

This repo is mainly for benchmarking several GC-friendly in-memory cache project.

Cache bench
===

## Benchmark target

1. [FastCache](https://github.com/VictoriaMetrics/fastcache)
2. [BigCache](https://github.com/allegro/bigcache)
3. [FreeCache](https://github.com/coocood/freecache)


## Benchmark result

**FastCache**

```
[root@dmj3 cache-bench]# bench-getset --cache fastcache --duration 10m --size 4294967296
**************** Benchmark Report ****************
[1] TPS 2009476
[2] MemStat Alloc: 2157MB, TotalAlloc:3781MB, Sys:2336MB, Mallocs:7559642, Free:1206375
[3] CacheStat map[]

[root@dmj3 cache-bench]# bench-getset --cache freecache --duration 10m --size 4294967296
```

**BigCache**

```shell
[root@dmj3 cache-bench]# bench-getset --cache bigcache --duration 10m --size 4294967296
**************** Benchmark Report ****************
[1] TPS 224461
[2] MemStat Alloc: 5053MB, TotalAlloc:21423MB, Sys:9616MB, Mallocs:266461416, Free:266452733
[3] CacheStat map[collisions:0 delHits:0 delMisses:0 hits:45584616 misses:21748802]
```

**FreeCache**

```
[root@dmj3 cache-bench]# bench-getset --cache freecache --duration 10m --size 4294967296


**************** Benchmark Report ****************
[1] TPS 995100
[2] MemStat Alloc: 5392MB, TotalAlloc:18154MB, Sys:9460MB, Mallocs:120023359, Free:113785828
[3] CacheStat map[avgAccess:1566268251 count:27531776 hit:120020207 hitRate:0.4020446295768596 miss:178504380 overWrite:0]
```

From the benchmark result, we can see that **fastCache** has best performance. In the mean time, it has lowest memory usage and memory allocation(which benefits from the off-heap allocation and chunk type memory management).

