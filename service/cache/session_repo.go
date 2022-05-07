package cache

import (
	"app"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

const (
	prefixUserBroker = "UB:" // user broker
	prefixSUID       = "SU:" // session user ids
)

type SessionRepository struct {
	Conn *redis.Client
	repo app.RepoSession
}

func NewSessionRepository(repo app.RepoSession) (*SessionRepository, error) {

	h, err := getHandle()
	if err != nil {
		return nil, err
	}

	return &SessionRepository{
		Conn: h.Conn,
		repo: repo,
	}, nil
}

// DelBrokerId 用户BROKER映射
func (sr *SessionRepository) DelBrokerId(ctx context.Context, userId int64) error {

	key := prefixUserBroker + strconv.FormatInt(userId, 10)
	cmd := sr.Conn.Del(ctx, key)
	return cmd.Err()
}

// SetBrokerId 用户BROKER映射
func (sr *SessionRepository) SetBrokerId(ctx context.Context, userId int64, brokerId int64) error {

	key := prefixUserBroker + strconv.FormatInt(userId, 10)
	cmd := sr.Conn.Set(ctx, key, brokerId, 0)
	return cmd.Err()
}

// GetBrokerId 用户BROKER映射
func (sr *SessionRepository) GetBrokerId(ctx context.Context, userId int64) (int64, error) {

	key := prefixUserBroker + strconv.FormatInt(userId, 10)
	cmd := sr.Conn.Get(ctx, key)
	return cmd.Int64()
}

// SetSessionUID 会话用户映射
func (sr *SessionRepository) SetSessionUID(ctx context.Context, sid int64, userIds []int64) error {

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

	go func() {
		if sr.repo.SetUID(sid, userIds) != nil {
			log.Println("error when write to database, SetSessionUID")
		}
	}()
	return cmd.Err()
}

// GetSessionUID 会话用户映射
func (sr *SessionRepository) GetSessionUID(ctx context.Context, sid int64) ([]int64, error) {

	key := prefixSUID + strconv.FormatInt(sid, 10)
	cmd := sr.Conn.Get(ctx, key)

	// get from repo
	if cmd.Err() != nil {
		ids, err := sr.repo.GetUID(sid)
		if err != nil {
			return nil, err
		}
		go func() {
			sr.SetSessionUID(ctx, sid, ids)
		}()
		return ids, err
	}

	s := strings.Split(cmd.Val(), ",")
	var n []int64
	for _, id := range s {
		i, _ := strconv.ParseInt(id, 10, 64)
		n = append(n, i)
	}
	return n, nil
}
