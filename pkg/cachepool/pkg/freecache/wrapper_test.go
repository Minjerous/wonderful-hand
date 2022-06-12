package freecache

import (
	"encoding/json"
	"errors"
	common "github.com/igxnon/cachepool/pkg/cache"
	"strconv"
	"sync"
	"testing"
	"time"
)

type Bar struct {
	Bar int64
	Yee string
}

type MyCoder struct {
}

func (m MyCoder) Encode(v interface{}) ([]byte, error) {
	if val, ok := v.(Bar); ok {
		return json.Marshal(val)
	}
	return nil, errors.New("not a Bar")
}

func (m MyCoder) Decode(b []byte) (interface{}, error) {
	var bar Bar
	err := json.Unmarshal(b, &bar)
	return bar, err
}

func TestCache(t *testing.T) {
	cache := New(time.Minute*5, MyCoder{}, 1024*1024)
	cache.SetDefault("foo", Bar{
		Bar: 0,
		Yee: "yee",
	})
	bar, ok := cache.Get("foo")
	if !ok {
		t.Error("not ok")
	}
	if bar.(Bar).Yee != "yee" {
		t.Error("not yee")
	}
}

func BenchmarkCacheGetExpiring(b *testing.B) {
	benchmarkCacheGet(b, 5*time.Minute)
}

func BenchmarkCacheGetNotExpiring(b *testing.B) {
	benchmarkCacheGet(b, common.NoExpiration)
}

func benchmarkCacheGet(b *testing.B, exp time.Duration) {
	b.StopTimer()
	tc := New(exp, MyCoder{}, 1024*1024)
	tc.Set("foobarba", "zquux", common.DefaultExpiration)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.Get("foobarba")
	}
}

func BenchmarkCacheGetManyConcurrentExpiring(b *testing.B) {
	benchmarkCacheGetManyConcurrent(b, 5*time.Minute)
}

func BenchmarkCacheGetManyConcurrentNotExpiring(b *testing.B) {
	benchmarkCacheGetManyConcurrent(b, common.NoExpiration)
}

func benchmarkCacheGetManyConcurrent(b *testing.B, exp time.Duration) {
	b.StopTimer()
	n := 10000
	tsc := New(exp, MyCoder{}, 1024*1024)
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
