package test

import (
	"context"
	"github.com/igxnon/cachepool"
	"github.com/igxnon/cachepool/helper"
	"github.com/igxnon/cachepool/pkg/go-cache"
	"github.com/streadway/amqp"
	"sync"
	"testing"
	"time"
)

// Test sync with MQ
func TestMQInOnePool(t *testing.T) {
	pool := cachepool.New()
	// use message queue, sync some cache
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ch, _ := conn.Channel()
	pool.UseMQ(context.Background(), ch, "cache")

	// publish an importance message into cache (
	_ = helper.Publish(ch, "下北沢一番臭の伝説", 114514, time.Minute*5)

	// sleep for a second
	time.Sleep(time.Second)

	got, exp, ok := pool.GetWithExpiration("下北沢一番臭の伝説")
	if ok {
		t.Log(got, exp)
	}

	time.Sleep(time.Second)

	// stop using message queue
	pool.StopMQ()
}

func TestMQInManyPool(t *testing.T) {
	var (
		p1      = cachepool.New()
		p2      = cachepool.New()
		p3      = cachepool.New()
		conn, _ = amqp.Dial("amqp://guest:guest@localhost:5672/")
		ch1, _  = conn.Channel()
		ch2, _  = conn.Channel()
		ch3, _  = conn.Channel()
	)

	p1.UseMQ(context.Background(), ch1, "cache1")
	p2.UseMQ(context.Background(), ch2, "cache2")
	p3.UseMQ(context.Background(), ch3, "cache3")

	// publish an importance message into cache (
	_ = helper.Publish(ch1, "下北沢一番臭の伝説", struct {
		Age    int
		Prefix string
		Movie  string
	}{
		24,
		"野獣せんべい",
		"真夏の夜の银夢",
	}, time.Minute*5)

	// sleep for a second
	time.Sleep(time.Second)

	got, exp, ok := p1.GetWithExpiration("下北沢一番臭の伝説")
	if ok {
		t.Log(got, exp)
	}

	got, exp, ok = p2.GetWithExpiration("下北沢一番臭の伝説")
	if ok {
		t.Log(got, exp)
	}

	got, exp, ok = p3.GetWithExpiration("下北沢一番臭の伝説")
	if ok {
		t.Log(got, exp)
	}

	time.Sleep(time.Second)

	// stop using message queue
	p1.StopMQ()
	p2.StopMQ()
	p3.StopMQ()
}

func BenchmarkCachePoolGet(b *testing.B) {
	b.StopTimer()
	pool := cachepool.New(cachepool.WithCache(gocache.NewCache(time.Minute*5, time.Minute*30)))
	pool.SetDefault("foo", "bar")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		pool.Get("foo")
	}
}

func BenchmarkSyncMapCachePoolGet(b *testing.B) {
	b.StopTimer()
	pool := cachepool.New(
		cachepool.WithCache(gocache.NewSyncMapCache(time.Minute*5, time.Minute*30)))
	pool.SetDefault("foo", "bar")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		pool.Get("foo")
	}
}

func BenchmarkCachePoolParallelGet(b *testing.B) {
	b.StopTimer()
	pool := cachepool.New(
		cachepool.WithCache(gocache.NewCache(time.Minute*5, time.Minute*30)))
	pool.SetDefault("foo", "bar")
	n := 100
	each := b.N / n
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			for j := 0; j < each; j++ {
				pool.Get("foo")
			}
			wg.Done()
		}()
	}
	b.StartTimer()
	wg.Wait()
}

func BenchmarkSyncMapCachePoolParallelGet(b *testing.B) {
	b.StopTimer()
	pool := cachepool.New(
		cachepool.WithCache(gocache.NewSyncMapCache(time.Minute*5, time.Minute*30)))
	pool.SetDefault("foo", "bar")
	n := 100
	each := b.N / n
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			for j := 0; j < each; j++ {
				pool.Get("foo")
			}
			wg.Done()
		}()
	}
	b.StartTimer()
	wg.Wait()
}
