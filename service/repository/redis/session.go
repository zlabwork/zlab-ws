package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"strings"
)

const (
	prefixOnline = "OL:"
	prefixSUID   = "SU:"
)

type SessionRepository struct {
	Conn *redis.Client
}

func NewOnlineRepository() (*SessionRepository, error) {

	h, err := getHandle()
	if err != nil {
		return nil, err
	}

	return &SessionRepository{
		Conn: h.Conn,
	}, nil
}

func (sr *SessionRepository) Close() error {

	return sr.Conn.Close()
}

// SetOffline 用户BROKER映射
func (sr *SessionRepository) SetOffline(ctx context.Context, userId int64) error {

	key := prefixOnline + strconv.FormatInt(userId, 10)
	cmd := sr.Conn.Del(ctx, key)
	return cmd.Err()
}

// SetOnline 用户BROKER映射
func (sr *SessionRepository) SetOnline(ctx context.Context, userId int64, node int64) error {

	key := prefixOnline + strconv.FormatInt(userId, 10)
	cmd := sr.Conn.Set(ctx, key, node, 0)
	return cmd.Err()
}

// GetOnline 用户BROKER映射
func (sr *SessionRepository) GetOnline(ctx context.Context, userId int64) (int64, error) {

	key := prefixOnline + strconv.FormatInt(userId, 10)
	cmd := sr.Conn.Get(ctx, key)
	return cmd.Int64()
}

// SetSUID 会话用户映射
func (sr *SessionRepository) SetSUID(ctx context.Context, sid int64, userIds []int64) error {

	if len(userIds) < 2 {
		return fmt.Errorf("error length of userIds")
	}

	var s []string
	for _, id := range userIds {
		u := strconv.FormatInt(id, 10)
		s = append(s, u)
	}
	key := prefixSUID + strconv.FormatInt(sid, 10)
	data := strings.Join(s, ",")
	cmd := sr.Conn.Set(ctx, key, data, 0)
	return cmd.Err()
}

// GetSUID 会话用户映射
func (sr *SessionRepository) GetSUID(ctx context.Context, sid int64) ([]int64, error) {

	key := prefixSUID + strconv.FormatInt(sid, 10)
	cmd := sr.Conn.Get(ctx, key)
	s := strings.Split(cmd.Val(), ",")

	var n []int64
	for _, id := range s {
		i, _ := strconv.ParseInt(id, 10, 64)
		n = append(n, i)
	}
	return n, nil
}
