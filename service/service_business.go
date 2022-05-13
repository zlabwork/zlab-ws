package service

import (
	"app"
	"app/service/business"
	"app/service/cache"
	"app/service/repository/mysql"
	"context"
	"encoding/binary"
	"fmt"
	"log"
)

type Business struct {
	message app.RepoMessage
	session app.CacheSession
	data    chan []byte

	// broker map
	bm map[int64]string
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
		data:    make(chan []byte),
		bm:      make(map[int64]string), // TODO:: 补充地址映射字典
	}, nil
}

func (bu *Business) Run() {

	go func() {
		business.Main(bu.data)
	}()

	go func() {
		for {
			select {
			case m := <-bu.data:

				sid := int64(binary.BigEndian.Uint64(m[4:12]))
				seq := int64(binary.BigEndian.Uint64(m[12:20]))
				send := int64(binary.BigEndian.Uint64(m[20:28]))
				msg := m[28:]
				fmt.Println("processing:", sid, seq, send, string(msg))

				bu.getBrokers(context.TODO(), sid)
			}
		}
	}()

}

// 根据 sid 获取 node 地址
func (bu *Business) getBrokers(ctx context.Context, sid int64) (map[int64]string, error) {

	userIds, err := bu.session.GetSessionUID(ctx, sid)
	if err != nil {
		return nil, err
	}

	m := make(map[int64]string)
	for _, uid := range userIds {

		// get broker node id
		n, err := bu.session.GetBrokerId(ctx, uid)
		if err != nil {
			log.Printf("no broker id for user %d\n", uid)
			continue
		}

		// node id to address
		addr, ok := bu.bm[n]
		if !ok {
			log.Println("broker node id is not in broker map")
			continue
		}
		m[uid] = addr
	}

	return m, nil
}
