package app

import (
	"context"
	"time"
)

type CacheToken interface {
	GetToken(id string) (string, error)
	SetToken(id string, token string) error
}

type CacheSession interface {
	DelBrokerId(ctx context.Context, userId int64) error
	SetBrokerId(ctx context.Context, userId int64, brokerId int64) error
	GetBrokerId(ctx context.Context, userId int64) (int64, error)
	SetSessionUID(ctx context.Context, sid int64, data []int64) error
	GetSessionUID(ctx context.Context, sid int64) ([]int64, error)
}

type RepoMessage interface {
	GetTodo(userId int64) ([]*RepoMsg, error)
	DeleteTodo(userId int64) error
	CreateTodo(msgType uint8, msgId string, sender, receiver int64, body []byte, date time.Time) error
	CreateLogs(msgType uint8, msgId string, sender, receiver int64, body []byte, date time.Time) error
}

type RepoSession interface {
	NewUID(sid int64, userIds []int64) error
	SetUID(sid int64, userIds []int64) error
	GetUID(sid int64) ([]int64, error)
}
