package redis

import (
	"github.com/go-redis/redis/v8"
)

type CacheRepository struct {
	Conn *redis.Client
}

func NewCacheRepository() (*CacheRepository, error) {

	h, err := getHandle()
	if err != nil {
		return nil, err
	}

	return &CacheRepository{
		Conn: h.Conn,
	}, nil
}

func (cr *CacheRepository) Close() error {

	return cr.Conn.Close()
}

func (cr *CacheRepository) GetToken(id string) (string, error) {

	return "", nil
}

func (cr *CacheRepository) SetToken(id string, token string) error {

	return nil
}
