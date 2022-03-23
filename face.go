package app

import "time"

type CacheFace interface {
	Close() error
	GetToken(id string) (string, error)
	SetToken(id string, token string) error
}

type RepoFace interface {
	GetTodo(userId int64) ([]*RepoMsg, error)
	DeleteTodo(userId int64) error
	CreateTodo(msgType uint8, msgId string, sender, receiver int64, body []byte, date time.Time) error
	CreateLogs(msgType uint8, msgId string, sender, receiver int64, body []byte, date time.Time) error
}
