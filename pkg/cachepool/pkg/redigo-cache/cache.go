package redicache

import (
	"github.com/gomodule/redigo/redis"
	common "github.com/igxnon/cachepool/pkg/cache"
	"strconv"
	"time"
)

var _ common.ICache = (*GlobalCache)(nil)

type GlobalCache struct {
	conn              redis.Conn
	defaultExpiration time.Duration
	coder             common.Coder
}

func (g *GlobalCache) set(k string, x interface{}, d time.Duration, norX string) error {
	b, err := g.coder.Encode(x)
	if err != nil {
		return err
	}
	if d == common.DefaultExpiration {
		d = g.defaultExpiration
	}
	if d > 0 {
		if norX == "" {
			_, err := g.conn.Do("SET", k, b, "PX",
				strconv.FormatInt(d.Milliseconds(), 10))
			return err
		}
		_, err := g.conn.Do("SET", k, b, "PX",
			strconv.FormatInt(d.Milliseconds(), 10), norX)
		return err
	}
	// no expire
	if norX == "" {
		_, err := g.conn.Do("SET", k, b)
		return err
	}
	_, err = g.conn.Do("SET", k, b, norX)
	return err
}

func (g *GlobalCache) Set(k string, x interface{}, d time.Duration) {
	_ = g.set(k, x, d, "")
}

func (g *GlobalCache) SetDefault(k string, x interface{}) {
	g.Set(k, x, common.DefaultExpiration)
}

// Add always return nil because redis keep adding once, if an error occurred
// while sending command to redis server, the error will be returned
func (g *GlobalCache) Add(k string, x interface{}, d time.Duration) error {
	return g.set(k, x, d, "NX")
}

func (g *GlobalCache) Replace(k string, x interface{}, d time.Duration) error {
	return g.set(k, x, d, "XX")
}

// Get return bytes, you should Unmarshal it in person
func (g *GlobalCache) Get(k string) (interface{}, bool) {
	b, err := redis.Bytes(g.conn.Do("GET", k))
	if err != nil {
		return nil, false
	}
	v, err := g.coder.Decode(b)
	return v, err == nil
}

func (g *GlobalCache) GetWithExpiration(k string) (interface{}, time.Time, bool) {
	ttl, err := redis.Int64(g.conn.Do("PTTL", k))
	if err != nil {
		return nil, time.Time{}, false
	}
	if ttl > 0 {
		exp := time.UnixMilli(ttl)
		b, ok := g.Get(k)
		return b, exp, ok
	}
	return nil, time.Time{}, false
}

func (g *GlobalCache) Increment(k string, n int64) error {
	_, err := g.conn.Do("INCRBY", k, n)
	return err
}

func (g *GlobalCache) Decrement(k string, n int64) error {
	_, err := g.conn.Do("DECRBY", k, n)
	return err
}

func (g *GlobalCache) Delete(k string) {
	_, _ = g.conn.Do("DEL", k)
}

func (g *GlobalCache) ItemCount() int {
	cnt, err := redis.Int(g.conn.Do("DBSIZE"))
	if err != nil {
		return -1
	}
	return cnt
}

func (g *GlobalCache) Flush() {
	// you'd better not do this
	return
}

func NewGlobalCache(defaultExpiration time.Duration, conn redis.Conn, coder common.Coder) *GlobalCache {
	return &GlobalCache{
		defaultExpiration: defaultExpiration,
		conn:              conn,
		coder:             coder,
	}
}
