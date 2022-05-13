package service

import (
	"app"
	"app/service/business"
	"app/service/cache"
	"app/service/repository/mysql"
	"encoding/binary"
	"fmt"
)

type Business struct {
	message app.RepoMessage
	session app.CacheSession
	data    chan []byte
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
			}
		}
	}()

}
