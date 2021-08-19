package mysql

import (
    "fmt"
    "log"
    "zlabws"
)

func getHandle() (*handle, error) {
    c := zlabws.Cfg.Db.Mysql
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci", c.User, c.Pass, c.Host, c.Port, c.Db)
    return ConnectMySQL(dsn)
}

func NewMessageService() *MessageService {
    h, err := getHandle()
    if err != nil {
        log.Println(err)
    }
    return &MessageService{h: h}
}
