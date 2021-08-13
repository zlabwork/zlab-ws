package redis

import (
    "github.com/go-redis/redis"
    "zlabws"
)

type handle struct {
    Conn *redis.Client
}

func ConnectRedis(c *zlabws.App) (*handle, error) {

    cli := redis.NewClient(&redis.Options{
        Addr: c.Database.Redis.Host + ":" + c.Database.Redis.Port,
        // Password: c.Pass,
        // DB:       c.Name,
    })
    return &handle{
        Conn: cli,
    }, nil
}
