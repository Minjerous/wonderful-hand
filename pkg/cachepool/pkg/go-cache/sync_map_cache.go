package gocache

import (
	"fmt"
	common "github.com/igxnon/cachepool/pkg/cache"
	"runtime"
	"sync"
	"time"
)

// Sync map Cache use golang sync.Map as its container
// sync.Map has fast query speed but nearly 0.5x slower
// than a map with mutex.
// Compared to sharded cache, cache based on sync.Map does
// not need to care about how to deal with sharded count
// growing while cache getting larger.

// NOTE: Some method maybe not atomic

var _ common.ICache = (*SyncMapCache)(nil)

type SyncMapCache struct {
	*syncMapCache
}

type syncMapCache struct {
	defaultExpiration time.Duration
	items             *sync.Map
	onEvicted         func(string, interface{})
	janitor           *janitor
	mu                sync.RWMutex
}

func (s *syncMapCache) Set(k string, x interface{}, d time.Duration) {
	// "Inlining" of set
	var e int64
	if d == common.DefaultExpiration {
		d = s.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	s.items.Store(k, Item{
		Object:     x,
		Expiration: e,
	})
}

func (s *syncMapCache) SetDefault(k string, x interface{}) {
	s.Set(k, x, common.DefaultExpiration)
}

// Add NOTE: 2x locking delay(s.mu and sync.Map set mutex)
func (s *syncMapCache) Add(k string, x interface{}, d time.Duration) error {
	s.mu.RLock() // avoid 2 Add method invoke at the same time and both return nil
	// Add method will change `ok` to true after it invokes
	_, ok := s.items.Load(k)
	if ok {
		s.mu.RUnlock()
		return fmt.Errorf("Item %s already exists", k)
	}
	s.Set(k, x, d)
	s.mu.RUnlock()
	return nil
}

func (s *syncMapCache) Replace(k string, x interface{}, d time.Duration) error {
	// we do not protect this because `Replace` does not change `ok`
	_, ok := s.items.Load(k)
	if !ok {
		return fmt.Errorf("Item %s doesn't exist", k)
	}
	s.Set(k, x, d)
	return nil
}

func (s *syncMapCache) Get(k string) (interface{}, bool) {
	item, ok := s.items.Load(k)
	if !ok {
		return nil, false
	}
	i := item.(Item)
	if i.Expiration > 0 {
		if time.Now().UnixNano() > i.Expiration {
			return nil, false
		}
	}
	return i.Object, true
}

func (s *syncMapCache) GetWithExpiration(k string) (interface{}, time.Time, bool) {
	item, ok := s.items.Load(k)
	if !ok {
		return nil, time.Time{}, false
	}

	i := item.(Item)
	if i.Expiration > 0 {
		if time.Now().UnixNano() > i.Expiration {
			return nil, time.Time{}, false
		}
		return i.Object, time.Unix(0, i.Expiration), true
	}

	return i.Object, time.Time{}, true
}

// Increment Cannot do Increment(k , n*-1) for uint
func (s *syncMapCache) Increment(k string, n int64) error {
	s.mu.RLock()
	item, ok := s.items.Load(k)
	if !ok || item.(Item).Expired() {
		s.mu.RUnlock()
		return fmt.Errorf("Item %s not found", k)
	}

	v := item.(Item)
	switch v.Object.(type) {
	case int:
		v.Object = v.Object.(int) + int(n)
	case int8:
		v.Object = v.Object.(int8) + int8(n)
	case int16:
		v.Object = v.Object.(int16) + int16(n)
	case int32:
		v.Object = v.Object.(int32) + int32(n)
	case int64:
		v.Object = v.Object.(int64) + n
	case uint:
		v.Object = v.Object.(uint) + uint(n)
	case uintptr:
		v.Object = v.Object.(uintptr) + uintptr(n)
	case uint8:
		v.Object = v.Object.(uint8) + uint8(n)
	case uint16:
		v.Object = v.Object.(uint16) + uint16(n)
	case uint32:
		v.Object = v.Object.(uint32) + uint32(n)
	case uint64:
		v.Object = v.Object.(uint64) + uint64(n)
	case float32:
		v.Object = v.Object.(float32) + float32(n)
	case float64:
		v.Object = v.Object.(float64) + float64(n)
	default:
		s.mu.RUnlock()
		return fmt.Errorf("The value for %s is not an integer", k)
	}
	s.items.Store(k, v)
	s.mu.RUnlock()
	return nil
}

func (s *syncMapCache) Decrement(k string, n int64) error {
	s.mu.RLock()
	item, ok := s.items.Load(k)
	if !ok || item.(Item).Expired() {
		s.mu.RUnlock()
		return fmt.Errorf("Item %s not found", k)
	}

	v := item.(Item)
	switch v.Object.(type) {
	case int:
		v.Object = v.Object.(int) - int(n)
	case int8:
		v.Object = v.Object.(int8) - int8(n)
	case int16:
		v.Object = v.Object.(int16) - int16(n)
	case int32:
		v.Object = v.Object.(int32) - int32(n)
	case int64:
		v.Object = v.Object.(int64) - n
	case uint:
		v.Object = v.Object.(uint) - uint(n)
	case uintptr:
		v.Object = v.Object.(uintptr) - uintptr(n)
	case uint8:
		v.Object = v.Object.(uint8) - uint8(n)
	case uint16:
		v.Object = v.Object.(uint16) - uint16(n)
	case uint32:
		v.Object = v.Object.(uint32) - uint32(n)
	case uint64:
		v.Object = v.Object.(uint64) - uint64(n)
	case float32:
		v.Object = v.Object.(float32) - float32(n)
	case float64:
		v.Object = v.Object.(float64) - float64(n)
	default:
		s.mu.RUnlock()
		return fmt.Errorf("The value for %s is not an integer", k)
	}
	s.items.Store(k, v)
	s.mu.RUnlock()
	return nil
}

func (s *syncMapCache) Delete(k string) {
	s.items.Delete(k)
}

func (s *syncMapCache) DeleteExpired() {

	now := time.Now().UnixNano()

	s.items.Range(func(k, item any) bool {
		if item.(Item).Expiration > 0 && now > item.(Item).Expiration {
			s.items.Delete(k)
		}
		return true
	})
}

func (s *syncMapCache) Items() map[string]common.IItem {
	items := make(map[string]common.IItem)
	s.items.Range(func(k, item any) bool {
		items[k.(string)] = item.(Item)
		return true
	})
	return items
}

func (s *syncMapCache) ItemCount() int {
	cnt := 0
	s.items.Range(func(_, _ any) bool {
		cnt++
		return true
	})
	return cnt
}

func (s *syncMapCache) Flush() {
	s.items = &sync.Map{}
}

func newSyncMapCache(de time.Duration) *syncMapCache {
	if de == 0 {
		de = -1
	}
	c := &syncMapCache{
		defaultExpiration: de,
		items:             &sync.Map{},
	}
	return c
}

func runSyncMapJanitor(c *syncMapCache, ci time.Duration) {
	j := &janitor{
		Interval: ci,
		stop:     make(chan bool),
	}
	c.janitor = j
	go j.Run(c)
}

func stopSyncMapJanitor(c *SyncMapCache) {
	c.janitor.stop <- true
}

func newSyncMapCacheWithJanitor(de time.Duration, ci time.Duration) *SyncMapCache {
	c := newSyncMapCache(de)
	// This trick ensures that the janitor goroutine (which--granted it
	// was enabled--is running DeleteExpired on c forever) does not keep
	// the returned C object from being garbage collected. When it is
	// garbage collected, the finalizer stops the janitor goroutine, after
	// which c can be collected.
	C := &SyncMapCache{c}
	if ci > 0 {
		runSyncMapJanitor(c, ci)
		runtime.SetFinalizer(C, stopSyncMapJanitor)
	}
	return C
}

func NewSyncMapCache(defaultExpiration, cleanupInterval time.Duration) *SyncMapCache {
	return newSyncMapCacheWithJanitor(defaultExpiration, cleanupInterval)
}
