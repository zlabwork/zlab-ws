package app

import "time"

type CacheFace interface {
	Close() error
	GetToken(id string) (string, error)
	SetToken(id string, token string) error
}

type RepoFace interface {
	Create(msgType uint8, msgId string, sender, receiver int64, body []byte, date time.Time) error
	CreateTodo(msgType uint8, msgId string, sender, receiver int64, body []byte, date time.Time) error
}
