package service

import (
	"app"
	"app/service/cache"
)

type Login struct {
	cache app.CacheToken
}

func NewLoginService() (*Login, error) {

	cache, err := cache.NewTokenRepository()
	if err != nil {
		return nil, err
	}

	return &Login{
		cache: cache,
	}, nil
}

func (lg *Login) CheckToken(token string) bool {

	return false
}
