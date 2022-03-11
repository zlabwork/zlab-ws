package redis

import (
	"github.com/go-redis/redis/v8"
)

type handle struct {
	Conn *redis.Client
}

// ConnectRedis
// redis://<user>:<pass>@localhost:6379/<db>
// https://redis.uptrace.dev/guide/server.html#connecting-to-redis-server
func ConnectRedis(dsn string) (*handle, error) {
	opt, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, err
	}
	cli := redis.NewClient(opt)
	return &handle{
		Conn: cli,
	}, nil
}
