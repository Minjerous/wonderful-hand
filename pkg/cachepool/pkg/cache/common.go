package cache

import "time"

const (
	// NoExpiration For use with functions that take an expiration time.
	NoExpiration time.Duration = -1
	// DefaultExpiration For use with functions that take an expiration time. Equivalent to
	// passing in the same expiration duration as was given to NewCache() or
	// NewCacheFrom() when the cache was created (e.g. 5 minutes.)
	DefaultExpiration time.Duration = 0
)

type ICache interface {
	// Set Add an item to the cache, replacing any existing item. If the duration is 0
	// (DefaultExpiration), the cache's default expiration time is used. If it is -1
	// (NoExpiration), the item never expires.
	Set(k string, x interface{}, d time.Duration)

	// SetDefault Add an item to the cache, replacing any existing item, using the default
	// expiration.
	SetDefault(k string, x interface{})

	// Add an item to the cache only if an item doesn't already exist for the given
	// key, or if the existing item has expired. Returns an error otherwise.
	Add(k string, x interface{}, d time.Duration) error

	// Replace Set a new value for the cache key only if it already exists, and the existing
	// item hasn't expired. Returns an error otherwise.
	Replace(k string, x interface{}, d time.Duration) error

	// Get an item from the cache. Returns the item or nil, and a bool indicating
	// whether the key was found.
	Get(k string) (interface{}, bool)

	// GetWithExpiration returns an item and its expiration time from the cache.
	// It returns the item or nil, the expiration time if one is set (if the item
	// never expires a zero value for time.Time is returned), and a bool indicating
	// whether the key was found.
	GetWithExpiration(k string) (interface{}, time.Time, bool)

	// Increment an item of type int, int8, int16, int32, int64, uintptr, uint,
	// uint8, uint32, or uint64, float32 or float64 by n. Returns an error if the
	// item's value is not an integer, if it was not found, or if it is not
	// possible to increment it by n. To retrieve the incremented value, use one
	// of the specialized methods, e.g. IncrementInt64.
	Increment(k string, n int64) error

	// Decrement an item of type int, int8, int16, int32, int64, uintptr, uint,
	// uint8, uint32, or uint64, float32 or float64 by n. Returns an error if the
	// item's value is not an integer, if it was not found, or if it is not
	// possible to decrement it by n. To retrieve the decremented value, use one
	// of the specialized methods, e.g. DecrementInt64.
	Decrement(k string, n int64) error

	// Delete an item from the cache. Does nothing if the key is not in the cache.
	Delete(k string)

	// ItemCount Returns the number of items in the cache. This may include items that have
	// expired, but have not yet been cleaned up.
	ItemCount() int

	// Flush Delete all items from the cache.
	Flush()
}

type IItem interface {
	// Expired Returns true if the item has expired.
	Expired() bool
}

// Coder for encoding some specified types and decode it,
// for all types supporting, you should use reflect
// fortunately, GlobalCacheSugar is a considerable way
type Coder interface {
	Encode(v interface{}) ([]byte, error)
	Decode(b []byte) (interface{}, error)
}
