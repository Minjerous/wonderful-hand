package cachepool

import (
	"database/sql"
	"github.com/igxnon/cachepool/pkg/cache"
	"time"
)

var _ ICachePool = (*DoubleCachePool)(nil)

// DoubleCachePool implement globalCache and act just like L1(localCache(readOnly map))
// L2(localCache) L3(globalCache) cache, and SQL database is just like Memory
// if globalCache implemented fits all type of the value it stored and Get() could
// return the value directly, helper.Query could be used on this pool
// DoubleCachePool used to build Cache-Aside between local cache and global cache
type DoubleCachePool struct {
	cache.ICache
	localCache  cache.ICache
	globalCache cache.ICache
	db          *sql.DB
}

func (c *DoubleCachePool) Set(k string, x interface{}, d time.Duration) {
	c.globalCache.Set(k, x, d)
	c.localCache.Delete(k)
	// todo publish delete
}

func (c *DoubleCachePool) SetDefault(k string, x interface{}) {
	c.Set(k, x, cache.DefaultExpiration)
}

func (c *DoubleCachePool) Add(k string, x interface{}, d time.Duration) error {
	err := c.globalCache.Add(k, x, d)
	if err != nil {
		return err
	}
	c.localCache.Delete(k)
	// todo publish delete
	return nil
}

func (c *DoubleCachePool) Replace(k string, x interface{}, d time.Duration) error {
	err := c.globalCache.Replace(k, x, d)
	if err != nil {
		return err
	}
	c.localCache.Delete(k)
	// todo publish delete
	return nil
}

func (c *DoubleCachePool) Get(k string) (interface{}, bool) {
	got, ok := c.localCache.Get(k)
	if ok {
		return got, ok
	}
	got, ok = c.globalCache.Get(k)
	if ok {
		c.localCache.SetDefault(k, got)
	}
	return got, ok
}

func (c *DoubleCachePool) GetWithExpiration(k string) (interface{}, time.Time, bool) {
	got, exp, ok := c.localCache.GetWithExpiration(k)
	if ok {
		return got, exp, ok
	}
	got, exp, ok = c.globalCache.GetWithExpiration(k)
	if ok {
		c.localCache.SetDefault(k, got)
	}
	return got, exp, ok
}

func (c *DoubleCachePool) Increment(k string, n int64) error {
	err := c.globalCache.Increment(k, n)
	if err != nil {
		return err
	}
	c.localCache.Delete(k)
	// todo publish delete
	return nil
}

func (c *DoubleCachePool) Decrement(k string, n int64) error {
	err := c.globalCache.Decrement(k, n)
	if err != nil {
		return err
	}
	c.localCache.Delete(k)
	// todo publish delete
	return nil
}

func (c *DoubleCachePool) Delete(k string) {
	c.localCache.Delete(k)
}

func (c *DoubleCachePool) ItemCount() int {
	return c.globalCache.ItemCount()
}

func (c *DoubleCachePool) Flush() {
	c.localCache.Flush()
}

func (c *DoubleCachePool) GetDatabase() *sql.DB {
	return c.db
}

func (c *DoubleCachePool) GetImplementedCache() cache.ICache {
	return c.ICache
}

func NewDouble(opt ...Option) *DoubleCachePool {
	opts := loadOptions(opt...)
	if opts._globalCache == nil {
		panic("global cache should be declared")
	}
	return &DoubleCachePool{
		ICache:      opts._globalCache,
		localCache:  opts.cache,
		globalCache: opts._globalCache,
		db:          opts.db,
	}
}
