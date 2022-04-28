package service

import (
	"app"
	"app/service/business"
	"app/service/cache"
	"app/service/repository/mysql"
)

type Business struct {
	message app.RepoMessage
	session app.CacheSession
}

func NewBusinessService() (*Business, error) {

	repo, err := mysql.NewSessionRepository()
	if err != nil {
		return nil, err
	}

	cache, err := cache.NewSessionRepository(repo)
	if err != nil {
		return nil, err
	}

	msg, err := mysql.NewMessageRepository()
	if err != nil {
		return nil, err
	}

	return &Business{
		message: msg,
		session: cache,
	}, nil
}

func (bu *Business) Run() {

	go func() {
		business.Main()
	}()
}
