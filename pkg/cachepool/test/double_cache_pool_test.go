package test

import (
	"github.com/gomodule/redigo/redis"
	"github.com/igxnon/cachepool"
	"github.com/igxnon/cachepool/pkg/cache"
	"github.com/igxnon/cachepool/pkg/go-cache"
	"testing"
	"time"
)

func TestDoubleCachePool(t *testing.T) {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		t.Error("redis does not connect")
	}

	pool := cachepool.NewDouble(
		cachepool.WithBuildinGlobalCache(time.Minute*30, conn, coder),
		cachepool.WithCache(gocache.NewCache(time.Minute*5, time.Minute*10)))

	pool.SetDefault("foo", Bar{Yee: "yee"})
	b, ok := pool.Get("foo")
	b, ok = pool.Get("foo")
	if !ok {
		t.Error("not ok")
	}

	var bar = b.(Bar)
	if bar.Yee != "yee" {
		t.Error("not yee")
	}
}

func BenchmarkDoubleCachePoolGet(b *testing.B) {
	benchmarkDoubleCachePoolGet(b, gocache.NewCache(time.Minute*5, time.Minute*10))
}

func BenchmarkDoubleSyncMapCachePoolGet(b *testing.B) {
	benchmarkDoubleCachePoolGet(b, gocache.NewSyncMapCache(time.Minute*5, time.Minute*10))
}

func benchmarkDoubleCachePoolGet(b *testing.B, c cache.ICache) {
	b.StopTimer()

	conn, _ := redis.Dial("tcp", "127.0.0.1:6379")

	pool := cachepool.NewDouble(
		cachepool.WithBuildinGlobalCache(time.Minute*30, conn, coder),
		cachepool.WithCache(c))

	pool.SetDefault("foo", Bar{Yee: "yee"})
	bytes, ok := pool.Get("foo")
	if !ok {
		b.Error("not ok")
	}

	var bar = bytes.(Bar)
	if bar.Yee != "yee" {
		b.Error("not yee")
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		pool.Get("foo")
	}
}
