package cache

import (
	"github.com/go-redis/redis/v8"
)

type TokenRepository struct {
	Conn *redis.Client
}

func NewTokenRepository() (*TokenRepository, error) {

	h, err := getHandle()
	if err != nil {
		return nil, err
	}

	return &TokenRepository{
		Conn: h.Conn,
	}, nil
}

func (tr *TokenRepository) GetToken(id string) (string, error) {

	return "", nil
}

func (tr *TokenRepository) SetToken(id string, token string) error {

	return nil
}
