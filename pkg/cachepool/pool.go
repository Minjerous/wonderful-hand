package cachepool

import (
	"context"
	"database/sql"
	"github.com/igxnon/cachepool/pkg/cache"
	"github.com/streadway/amqp"
	_ "unsafe"
)

var _ ICachePool = (*CachePool)(nil)

type ICachePool interface {
	cache.ICache

	GetDatabase() *sql.DB
	GetImplementedCache() cache.ICache
}

type CachePool struct {
	cache.ICache
	db       *sql.DB
	cancelMQ context.CancelFunc
}

func (c *CachePool) GetDatabase() *sql.DB {
	return c.db
}

func (c *CachePool) GetImplementedCache() cache.ICache {
	return c.ICache
}

// UseMQ uses rabbitmq to sync some cache between different machines.
// It returns a channel, if err happened before run mq listener, the error
// will be sent into the channel immediately. And after ctx done(StopMQ())
// or mq closed, nil will be sent into the channel.
// name passed to it must be a unique id among all machines
// it is useless for global cache such as NoSQL based cache
func (c *CachePool) UseMQ(ctx context.Context, ch *amqp.Channel, name string) <-chan error {
	ctx, cancel := context.WithCancel(ctx)
	c.cancelMQ = cancel
	cha := make(chan error)
	go func() {
		cha <- run_mq(ctx, c, ch, name)
	}()
	return cha
}

// StopMQ stop using message queue
func (c *CachePool) StopMQ() {
	if c.cancelMQ != nil {
		c.cancelMQ()
	}
}

//go:linkname run_mq github.com/igxnon/cachepool/helper.runSyncFromMQ
//noinspection ALL
func run_mq(ctx context.Context, cache cache.ICache, ch *amqp.Channel, name string) error

func New(opt ...Option) *CachePool {
	opts := loadOptions(opt...)
	return &CachePool{
		ICache: opts.cache,
		db:     opts.db,
	}
}
