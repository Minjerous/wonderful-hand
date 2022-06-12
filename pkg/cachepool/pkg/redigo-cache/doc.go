// Package redicache is some implements of ICache called GlobalCache,
// they use "github.com/gomodule/redigo/redis" as redis client and base
// on redis to build cache.
// GlobalCacheSugar provide methods that suit for all value types, but it
// cost performance to serialize value by reflect.
// GlobalCache initialized with a Coder that help encode and decode value,
// pass your own Coder if you want for performance.
package redicache
