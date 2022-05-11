package mysql

import (
	"app"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var _handle *handle

type handle struct {
	Conn *sql.DB
}

var (
	ErrNoRows        = sql.ErrNoRows
	DefaultCharset   = "utf8mb4"
	DefaultCollation = "utf8mb4_unicode_ci"
)

// ConnectMySQL
// @docs http://go-database-sql.org/retrieving.html
// @docs https://github.com/go-sql-driver/mysql/wiki/Examples
// username:password@tcp(localhost:3306)/dbname?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci
func ConnectMySQL(dsn string) (*handle, error) {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &handle{
		Conn: conn,
	}, nil
}

func (db *handle) CreateDatabase(database string) error {
	q := fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET = %s COLLATE = %s;", database, DefaultCharset, DefaultCollation)
	_, err := db.Conn.Exec(q)
	return err
}

func (db *handle) Drop(database string) error {
	q := fmt.Sprintf("DROP DATABASE %s;", database)
	_, err := db.Conn.Exec(q)
	return err
}

func getHandle() (*handle, error) {

	if _handle != nil {
		return _handle, nil
	}

	c := app.Yaml.Business.Db.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci", c.User, c.Pass, c.Host, c.Port, c.Name)

	var err error
	_handle, err = ConnectMySQL(dsn)
	if err != nil {
		return nil, err
	}
	return _handle, nil
}
