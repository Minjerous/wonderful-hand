package test

import (
	"database/sql"
	"encoding/json"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/igxnon/cachepool"
	"github.com/igxnon/cachepool/helper"
	"github.com/igxnon/cachepool/pkg/cache"
	"github.com/igxnon/cachepool/pkg/redigo-cache"
	"strconv"
	"sync"
	"testing"
	"time"
)

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

type Bar struct {
	Bar int64
	Yee string
	Foo sql.NullTime
}

var (
	conn, _ = redis.Dial("tcp", "127.0.0.1:6379")
	coder   = MyCoder{}
)

func TestGlobalCache(t *testing.T) {
	tc := redicache.NewGlobalCache(time.Minute*10, conn, coder)
	tc.Set("foobarba", Bar{Yee: "hello"}, cache.DefaultExpiration)
	g, ok := tc.Get("foobarba")
	var got = g.(Bar)
	if !ok || got.Yee != "hello" {
		t.Error("error")
	}
}

func TestGlobalCacheHelper(t *testing.T) {
	dsn := "root:12345678@tcp(127.0.0.1:3306)/awesome?charset=utf8mb4&parseTime=True"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	tc := redicache.NewGlobalCache(time.Minute*10, conn, coder)
	pool := cachepool.New(cachepool.WithCache(tc), cachepool.WithDatabase(db))
	got, err := helper.QueryRow[Bar](pool, "bar:combine", "SELECT * FROM t LIMIT 1 OFFSET 1")
	if err != nil {
		t.Error(err)
	}
	got, err = helper.QueryRow[Bar](pool, "bar:combine", "SELECT * FROM t LIMIT 1 OFFSET 1")
	t.Log(got)
}

func BenchmarkGlobalCacheGetExpiring(b *testing.B) {
	benchmarkGlobalCacheGet(b, 5*time.Minute)
}

func BenchmarkGlobalCacheGetNotExpiring(b *testing.B) {
	benchmarkGlobalCacheGet(b, cache.NoExpiration)
}

func benchmarkGlobalCacheGet(b *testing.B, exp time.Duration) {
	b.StopTimer()
	tc := redicache.NewGlobalCache(exp, conn, coder)
	tc.Set("foobarba", Bar{Yee: "hello"}, cache.DefaultExpiration)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.Get("foobarba")
	}
}

func BenchmarkGlobalCacheGetManyConcurrentExpiring(b *testing.B) {
	benchmarkGlobalCacheGetManyConcurrent(b, 5*time.Minute)
}

func BenchmarkGlobalCacheGetManyConcurrentNotExpiring(b *testing.B) {
	benchmarkGlobalCacheGetManyConcurrent(b, cache.NoExpiration)
}

func benchmarkGlobalCacheGetManyConcurrent(b *testing.B, exp time.Duration) {
	b.StopTimer()
	n := 10000
	tsc := redicache.NewGlobalCache(exp, conn, coder)
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		k := "foo" + strconv.Itoa(i)
		keys[i] = k
		tsc.Set(k, Bar{Yee: "hello"}, cache.DefaultExpiration)
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
