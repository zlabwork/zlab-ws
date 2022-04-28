package redis

import (
	"app"
	"fmt"
)

func getHandle() (*handle, error) {

	host := app.Yaml.Db.Redis.Host
	port := app.Yaml.Db.Redis.Port
	name := "0"
	dsn := fmt.Sprintf("redis://%s:%d/%s", host, port, name)
	return ConnectRedis(dsn)
}
