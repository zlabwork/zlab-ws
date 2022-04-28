package mysql

import (
	"app"
	"fmt"
)

var _handle *handle

func getHandle() (*handle, error) {

	if _handle != nil {
		return _handle, nil
	}

	c := app.Yaml.Db.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci", c.User, c.Pass, c.Host, c.Port, c.Name)

	var err error
	_handle, err = ConnectMySQL(dsn)
	if err != nil {
		return nil, err
	}
	return _handle, nil
}
