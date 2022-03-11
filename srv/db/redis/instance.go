package redis

import (
	"github.com/go-redis/redis/v8"
	"zlabws"
)

func NewRedisService() (*handle, error) {
	c := zlabws.Cfg.Db.Redis
	cli := redis.NewClient(&redis.Options{
		Addr: c.Host + ":" + c.Port,
	})
	return ConnectRedis(cli)
}
