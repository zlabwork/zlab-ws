package mysql

import (
	"fmt"
	"os"
)

func getHandle() (*handle, error) {
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASS")
	name := os.Getenv("MYSQL_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci", user, pass, host, port, name)
	return ConnectMySQL(dsn)
}

func NewMessageService() (*MessageService, error) {
	h, err := getHandle()
	if err != nil {
		return nil, err
	}
	return &MessageService{h: h}, nil
}
