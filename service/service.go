package service

import (
	"app"
	"app/service/redis"
	"app/service/repository/mysql"
	"app/service/ws"
	"net/http"
)

type Container struct {
	cache app.CacheFace
	repo  app.RepoFace
	hub   *ws.Hub
}

func NewService() (*Container, error) {

	cs, err := redis.NewCacheRepository()
	if err != nil {
		return nil, err
	}
	repo, err := mysql.NewRepoFace()
	if err != nil {
		return nil, err
	}
	hub := ws.NewHub(repo)

	return &Container{
		cache: cs,
		repo:  repo,
		hub:   hub,
	}, nil
}

func (co *Container) Run() {
	go co.hub.Run()
}

func (co *Container) ServeWs(w http.ResponseWriter, r *http.Request) {
	ws.ServeWs(co.hub, co.cache, co.repo, w, r)
}
