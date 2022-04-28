package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type SeqRepository struct {
	Conn *redis.Client
}

func NewSeqRepository() (*SeqRepository, error) {

	h, err := getHandle()
	if err != nil {
		return nil, err
	}

	return &SeqRepository{
		Conn: h.Conn,
	}, nil
}

func (sr *SeqRepository) GetSeqId(ctx context.Context, userId int64) (uint64, error) {
	k := "SEQ:" + strconv.FormatInt(userId, 10)
	cmd := sr.Conn.Incr(ctx, k)
	return cmd.Uint64()
}

func (sr *SeqRepository) SetSeqId(ctx context.Context, userId int64, value uint64) error {
	k := "SEQ:" + strconv.FormatInt(userId, 10)
	cmd := sr.Conn.Set(ctx, k, value, 0)
	return cmd.Err()
}

func (sr *SeqRepository) GetMaxSeq(ctx context.Context, section int64) (uint64, error) {
	k := "SEQS:" + strconv.FormatInt(section, 10)
	cmd := sr.Conn.Get(ctx, k)
	return cmd.Uint64()
}

func (sr *SeqRepository) SetMaxSeq(ctx context.Context, section int64, value uint64) error {
	k := "SEQS:" + strconv.FormatInt(section, 10)
	cmd := sr.Conn.Set(ctx, k, value, 0)
	return cmd.Err()
}

func (sr *SeqRepository) LoadMaxSeq(ctx context.Context, data map[int64]uint64) error {

	var m = make(map[string]interface{})
	for k, v := range data {
		key := "SEQS:" + strconv.FormatInt(k, 10)
		m[key] = v
	}
	cmd := sr.Conn.MSet(ctx, m)
	return cmd.Err()
}
