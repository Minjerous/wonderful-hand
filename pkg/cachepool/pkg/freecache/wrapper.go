package freecache

import (
	"fmt"
	internal "github.com/coocood/freecache"
	common "github.com/igxnon/cachepool/pkg/cache"
	"sync"
	"time"
)

// Cache wrap internal.Cache and implement ICache
type Cache struct {
	*internal.Cache
	coder             common.Coder
	defaultExpiration time.Duration
	mu                sync.RWMutex
}

func (c *Cache) set(k string, x interface{}, d time.Duration) error {
	if d == common.DefaultExpiration {
		d = c.defaultExpiration
	}
	b, err := c.coder.Encode(x)
	if err != nil {
		return err
	}
	return c.Cache.Set([]byte(k), b, int(d.Seconds()))
}

func (c *Cache) Set(k string, x interface{}, d time.Duration) {
	_ = c.set(k, x, d)
}

func (c *Cache) SetDefault(k string, x interface{}) {
	c.Set(k, x, common.DefaultExpiration)
}

func (c *Cache) Add(k string, x interface{}, d time.Duration) error {
	c.mu.RLock()
	_, err := c.Cache.Get([]byte(k))
	if err == nil {
		c.mu.RUnlock()
		return fmt.Errorf("Item %s already exists", k)
	}
	err = c.set(k, x, d)
	c.mu.RUnlock()
	return err
}

func (c *Cache) Replace(k string, x interface{}, d time.Duration) error {
	_, err := c.Cache.Get([]byte(k))
	if err != nil {
		return fmt.Errorf("Item %s is not exists", k)
	}
	return c.set(k, x, d)
}

func (c *Cache) Get(k string) (interface{}, bool) {
	b, err := c.Cache.Get([]byte(k))
	if err != nil {
		return nil, false
	}
	v, err := c.coder.Decode(b)
	return v, err == nil
}

func (c *Cache) GetWithExpiration(k string) (interface{}, time.Time, bool) {
	b, expireAt, err := c.Cache.GetWithExpiration([]byte(k))
	if err != nil {
		return nil, time.Time{}, false
	}
	v, err := c.coder.Decode(b)
	return v, time.Unix(int64(expireAt), 0), err == nil
}

func (c *Cache) Increment(k string, n int64) error {
	c.mu.RLock()
	v, ok := c.Get(k)
	if !ok {
		c.mu.RUnlock()
		return fmt.Errorf("Item %s is not exists", k)
	}
	switch v.(type) {
	case int:
		v = v.(int) + int(n)
	case int8:
		v = v.(int8) + int8(n)
	case int16:
		v = v.(int16) + int16(n)
	case int32:
		v = v.(int32) + int32(n)
	case int64:
		v = v.(int64) + n
	case uint:
		v = v.(uint) + uint(n)
	case uintptr:
		v = v.(uintptr) + uintptr(n)
	case uint8:
		v = v.(uint8) + uint8(n)
	case uint16:
		v = v.(uint16) + uint16(n)
	case uint32:
		v = v.(uint32) + uint32(n)
	case uint64:
		v = v.(uint64) + uint64(n)
	case float32:
		v = v.(float32) + float32(n)
	case float64:
		v = v.(float64) + float64(n)
	default:
		c.mu.RUnlock()
		return fmt.Errorf("The value for %s is not an integer", k)
	}
	err := c.set(k, v, c.defaultExpiration)
	c.mu.RUnlock()
	return err
}

func (c *Cache) Decrement(k string, n int64) error {
	c.mu.RLock()
	v, ok := c.Get(k)
	if !ok {
		c.mu.RUnlock()
		return fmt.Errorf("Item %s is not exists", k)
	}
	switch v.(type) {
	case int:
		v = v.(int) - int(n)
	case int8:
		v = v.(int8) - int8(n)
	case int16:
		v = v.(int16) - int16(n)
	case int32:
		v = v.(int32) - int32(n)
	case int64:
		v = v.(int64) - n
	case uint:
		v = v.(uint) - uint(n)
	case uintptr:
		v = v.(uintptr) - uintptr(n)
	case uint8:
		v = v.(uint8) - uint8(n)
	case uint16:
		v = v.(uint16) - uint16(n)
	case uint32:
		v = v.(uint32) - uint32(n)
	case uint64:
		v = v.(uint64) - uint64(n)
	case float32:
		v = v.(float32) - float32(n)
	case float64:
		v = v.(float64) - float64(n)
	default:
		c.mu.RUnlock()
		return fmt.Errorf("The value for %s is not an integer", k)
	}
	err := c.set(k, v, c.defaultExpiration)
	c.mu.RUnlock()
	return err
}

func (c *Cache) Delete(k string) {
	c.Cache.Del([]byte(k))
}

func (c *Cache) ItemCount() int {
	return int(c.Cache.EntryCount())
}

func (c *Cache) Flush() {
	c.Cache.Clear()
}

func New(defaultExpiration time.Duration, coder common.Coder, size int) *Cache {
	return &Cache{
		Cache:             internal.NewCache(size),
		coder:             coder,
		defaultExpiration: defaultExpiration,
		mu:                sync.RWMutex{},
	}
}
