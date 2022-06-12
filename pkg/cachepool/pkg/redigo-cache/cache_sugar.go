package redicache

import (
	"github.com/alecthomas/binary"
	"github.com/gomodule/redigo/redis"
	common "github.com/igxnon/cachepool/pkg/cache"
	"strconv"
	"time"
)

var _ common.ICache = (*GlobalCacheSugar)(nil)

type GlobalCacheSugar struct {
	conn              redis.Conn
	defaultExpiration time.Duration
}

func (g *GlobalCacheSugar) set(k string, x interface{}, d time.Duration, norX string) error {
	b, err := binary.Marshal(x)
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

func (g *GlobalCacheSugar) Set(k string, x interface{}, d time.Duration) {
	_ = g.set(k, x, d, "")
}

func (g *GlobalCacheSugar) SetDefault(k string, x interface{}) {
	g.Set(k, x, common.DefaultExpiration)
}

// Add always return nil because redis keep adding once, if an error occurred
// while sending command to redis server, the error will be returned
func (g *GlobalCacheSugar) Add(k string, x interface{}, d time.Duration) error {
	return g.set(k, x, d, "NX")
}

func (g *GlobalCacheSugar) Replace(k string, x interface{}, d time.Duration) error {
	return g.set(k, x, d, "XX")
}

// Get return bytes, you should Unmarshal it in person
func (g *GlobalCacheSugar) Get(k string) (interface{}, bool) {
	b, err := redis.Bytes(g.conn.Do("GET", k))
	if err != nil {
		return nil, false
	}
	return b, true
}

// GetUnmarshal helps unmarshal object, obj argument should be a pointer
// it implements interface unmarshalable in helper/internal.query
func (g *GlobalCacheSugar) GetUnmarshal(k string, obj interface{}) bool {
	b, ok := g.Get(k)
	if !ok {
		return false
	}
	return binary.Unmarshal(b.([]byte), obj) == nil
}

func (g *GlobalCacheSugar) GetWithExpiration(k string) (interface{}, time.Time, bool) {
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

func (g *GlobalCacheSugar) Increment(k string, n int64) error {
	_, err := g.conn.Do("INCRBY", k, n)
	return err
}

func (g *GlobalCacheSugar) Decrement(k string, n int64) error {
	_, err := g.conn.Do("DECRBY", k, n)
	return err
}

func (g *GlobalCacheSugar) Delete(k string) {
	_, _ = g.conn.Do("DEL", k)
}

func (g *GlobalCacheSugar) ItemCount() int {
	cnt, err := redis.Int(g.conn.Do("DBSIZE"))
	if err != nil {
		return -1
	}
	return cnt
}

func (g *GlobalCacheSugar) Flush() {
	// you'd better not do this
	return
}

func NewGlobalCacheSugar(defaultExpiration time.Duration, conn redis.Conn) *GlobalCacheSugar {
	return &GlobalCacheSugar{
		defaultExpiration: defaultExpiration,
		conn:              conn,
	}
}
