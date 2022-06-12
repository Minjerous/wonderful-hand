package gocache

import (
	common "github.com/igxnon/cachepool/pkg/cache"
	"strconv"
	"sync"
	"testing"
	"time"
)

func BenchmarkSyncMapCacheGetExpiring(b *testing.B) {
	benchmarkSyncMapCacheGet(b, 5*time.Minute)
}

func BenchmarkSyncMapCacheGetNotExpiring(b *testing.B) {
	benchmarkSyncMapCacheGet(b, common.NoExpiration)
}

func benchmarkSyncMapCacheGet(b *testing.B, exp time.Duration) {
	b.StopTimer()
	tc := NewSyncMapCache(exp, 0)
	tc.Set("foobarba", "zquux", common.DefaultExpiration)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.Get("foobarba")
	}
}

func BenchmarkSyncMapCacheGetManyConcurrentExpiring(b *testing.B) {
	benchmarkSyncMapCacheGetManyConcurrent(b, 5*time.Minute)
}

func BenchmarkSyncMapCacheGetManyConcurrentNotExpiring(b *testing.B) {
	benchmarkSyncMapCacheGetManyConcurrent(b, common.NoExpiration)
}

func benchmarkSyncMapCacheGetManyConcurrent(b *testing.B, exp time.Duration) {
	b.StopTimer()
	n := 10000
	tsc := NewSyncMapCache(exp, 0)
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		k := "foo" + strconv.Itoa(i)
		keys[i] = k
		tsc.Set(k, "bar", common.DefaultExpiration)
	}
	each := b.N / n
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for _, v := range keys {
		go func(k string) {
			for j := 0; j < each; j++ {
				tsc.Get(k)
			}
			wg.Done()
		}(v)
	}
	b.StartTimer()
	wg.Wait()
}
