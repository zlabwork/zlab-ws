package redis

import (
	"fmt"
	"os"
)

func getHandle() (*handle, error) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	name := "0"
	dsn := fmt.Sprintf("redis://%s:%s/%s", host, port, name)
	return ConnectRedis(dsn)
}
