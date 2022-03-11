package service

import (
	"app/service/redis"
	"net/http"
)

type CacheFace interface {
	Close() error
	GetToken(id string) (*string, error)
	SetToken(id string, token *string) error
}

type RepoFace interface {
}

type Container struct {
	Cache CacheFace
	Repo  RepoFace
	Hub   *Hub
}

func NewService() (*Container, error) {

	hub := newHub()
	cs, err := redis.NewCacheRepository()
	if err != nil {
		return nil, err
	}

	// TODO: repo
	return &Container{
		Cache: cs,
		Hub:   hub,
	}, nil
}

func (co *Container) Run() {
	go co.Hub.Run()
}

func (co *Container) ServeWs(w http.ResponseWriter, r *http.Request) {
	serveWs(co, w, r)
}
