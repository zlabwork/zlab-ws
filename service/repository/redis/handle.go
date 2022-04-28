package redis

import (
	"github.com/go-redis/redis/v8"
)

type handle struct {
	Conn *redis.Client
}

// ConnectRedis
// By default, the pool size is 10 connections per every available CPU as reported by runtime.GOMAXPROCS
// redis://<user>:<pass>@localhost:6379/<db>
// https://redis.uptrace.dev/guide/server.html#connecting-to-redis-server
// https://redis.uptrace.dev/guide/go-redis-debugging.html#connection-pool-size
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

// Cluster
// https://redis.uptrace.dev/guide/go-redis-cluster.html
//cli := redis.NewClusterClient(&redis.ClusterOptions{
//	Addrs: []string{":7000", ":7001", ":7002"},
//})

// Cache
// https://redis.uptrace.dev/guide/go-redis-cache.html
