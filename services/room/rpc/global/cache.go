package global

import (
	"github.com/gomodule/redigo/redis"
	"github.com/igxnon/cachepool"
	"github.com/igxnon/cachepool/pkg/freecache"
	"github.com/igxnon/cachepool/pkg/redigo-cache"
	"log"
	"time"
	"wonderful-hand-room/rpc/internal/config"
)

var CachePool cachepool.ICachePool

type coder struct {
}

func (c *coder) Encode(v interface{}) ([]byte, error) {
	// TODO 序列化 room
	panic("implement me")
}

func (c *coder) Decode(b []byte) (interface{}, error) {
	// TODO 反序列化 room
	panic("implement me")
}

func init() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalln(err)
	}
	conn, err := redis.Dial("tcp", cfg.Redis[0].Host,
		redis.DialPassword(cfg.Redis[0].Password))
	if err != nil {
		log.Fatalln(err)
	}
	CachePool = cachepool.NewDouble(
		cachepool.WithCache(freecache.New(time.Minute*20, &coder{}, 1024*1024*100)), // 100MB 本地cache
		cachepool.WithGlobalCache(redicache.NewGlobalCache(time.Minute*20, conn, &coder{})),
	)
}
