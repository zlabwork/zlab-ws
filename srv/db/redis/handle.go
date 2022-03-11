package redis

import (
	"github.com/go-redis/redis/v8"
)

type handle struct {
	Conn *redis.Client
}

func ConnectRedis(cli *redis.Client) (*handle, error) {
	return &handle{
		Conn: cli,
	}, nil
}
