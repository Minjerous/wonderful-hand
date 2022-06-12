package cachepool

import (
	"database/sql"
	"github.com/gomodule/redigo/redis"
	"github.com/igxnon/cachepool/pkg/cache"
	"github.com/igxnon/cachepool/pkg/go-cache"
	"github.com/igxnon/cachepool/pkg/redigo-cache"
	"time"
)

type Option func(*Options)

type Options struct {
	db           *sql.DB
	cache        cache.ICache
	_globalCache cache.ICache
}

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	if opts.cache == nil {
		opts.cache = gocache.NewCache(time.Minute*5, time.Minute*30)
	}
	return opts
}

func WithCache(cache cache.ICache) Option {
	return func(opt *Options) {
		opt.cache = cache
	}
}

func WithGlobalCache(cache cache.ICache) Option {
	return func(opt *Options) {
		opt._globalCache = cache
	}
}

// WithBuildinGlobalCache use buildin redis global cache
func WithBuildinGlobalCache(defaultExpiration time.Duration, conn redis.Conn, coder cache.Coder) Option {
	return func(opt *Options) {
		opt._globalCache = redicache.NewGlobalCache(defaultExpiration, conn, coder)
	}
}

func WithBuildinGlobalCacheSugar(defaultExpiration time.Duration, conn redis.Conn) Option {
	return func(opt *Options) {
		opt._globalCache = redicache.NewGlobalCacheSugar(defaultExpiration, conn)
	}
}

func WithDatabase(db *sql.DB) Option {
	return func(opt *Options) {
		opt.db = db
	}
}
