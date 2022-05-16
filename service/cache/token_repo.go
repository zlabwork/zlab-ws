package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	prefixTK = "TK:" // for token
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

func (tk *TokenRepository) GetToken(ctx context.Context, id string) (string, error) {

	cmd := tk.Conn.Get(ctx, prefixTK+id)
	if cmd.Err() != nil {
		return "", cmd.Err()
	}

	return cmd.String(), nil
}

func (tk *TokenRepository) SetToken(ctx context.Context, id string, token string) error {

	cmd := tk.Conn.Set(ctx, prefixTK+id, token, 24*time.Hour)
	return cmd.Err()
}
