package mysql

import (
	"fmt"
	"zlabws"
)

func getHandle() (*handle, error) {
	c := zlabws.Cfg.Db.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci", c.User, c.Pass, c.Host, c.Port, c.Db)
	return ConnectMySQL(dsn)
}

func NewMessageService() (*MessageService, error) {
	h, err := getHandle()
	if err != nil {
		return nil, err
	}
	return &MessageService{h: h}, nil
}
