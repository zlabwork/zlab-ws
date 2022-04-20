package service

import (
	"app"
	"app/service/business"
	"app/service/redis"
	"app/service/repository/mysql"
)

type Business struct {
	cache app.CacheFace
	repo  app.RepoFace
}

func NewBusinessService() (*Business, error) {

	cs, err := redis.NewCacheRepository()
	if err != nil {
		return nil, err
	}

	repo, err := mysql.NewRepoFace()
	if err != nil {
		return nil, err
	}

	return &Business{
		cache: cs,
		repo:  repo,
	}, nil
}

func (bu *Business) Run() {

	go func() {
		business.Main()
	}()
}
