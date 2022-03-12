package service

import (
	"app"
	"app/service/redis"
	"app/service/ws"
	"net/http"
)

type Container struct {
	Cache app.CacheFace
	Repo  app.RepoFace
	Hub   *ws.Hub
}

func NewService() (*Container, error) {

	hub := ws.NewHub()
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
	ws.ServeWs(co.Hub, co.Cache, w, r)
}
