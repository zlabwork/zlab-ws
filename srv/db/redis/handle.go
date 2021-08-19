package redis

import (
    "github.com/go-redis/redis"
)

type handle struct {
    Conn *redis.Client
}

func ConnectRedis(cli *redis.Client) (*handle, error) {
    return &handle{
        Conn: cli,
    }, nil
}
