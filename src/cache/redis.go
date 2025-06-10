package cache

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"muyu.com/reader_server_go/v1/src/config"
	"muyu.com/reader_server_go/v1/src/logger"
	"muyu.com/reader_server_go/v1/src/util"
	"strings"
)

var rdb *redis.Client

func Redis() *redis.Client {
	if rdb == nil {
		logger.Panic("redis has not been initialized!")
	}
	return rdb
}

func InitRedis() {
	con := config.Config.Redis

	host := strings.TrimSpace(con.Host)
	if host == "" {
		logger.Panic("redis host is blank!")
	}

	port := con.Port
	if port == 0 {
		port = 6379
	}

	poolSize := con.PoolSize
	poolSize = util.MaxInt(poolSize, 10)

	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		DB:       0, // use default
		PoolSize: poolSize,
	})
}
