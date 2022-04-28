package redis

import (
	"app"
	"fmt"
)

var _handle *handle

func getHandle() (*handle, error) {

	if _handle != nil {
		return _handle, nil
	}

	c := app.Yaml.Db.Redis
	name := "1"
	dsn := fmt.Sprintf("redis://%s:%d/%s", c.Host, c.Port, name)

	var err error
	_handle, err = ConnectRedis(dsn)
	if err != nil {
		return nil, err
	}
	return _handle, nil
}
