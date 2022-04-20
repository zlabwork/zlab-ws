package service

import (
	"app"
	"app/restful"
	"app/service/redis"
	"app/service/repository/mysql"
	"app/service/ws"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Container struct {
	cache app.CacheFace
	repo  app.RepoFace
	hub   *ws.Hub
}

func NewBrokerService() (*Container, error) {

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

func (co *Container) information() {

	type info struct {
		Node    string `json:"node"`
		Time    int64  `json:"time"`
		Clients int    `json:"clients"`
	}

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C

		bs, err := json.Marshal(&info{Node: os.Getenv("APP_NODE"), Time: time.Now().UTC().Unix(), Clients: co.hub.GetClientsNumber()})
		if err != nil {
			return
		}
		log.Println("INFO:", string(bs))
	}
}

func (co *Container) serveWs(w http.ResponseWriter, r *http.Request) {
	ws.ServeWs(co.hub, co.cache, co.repo, w, r)
}

func (co *Container) Run(addr *string) {

	go co.hub.Run()
	go co.information()

	http.HandleFunc("/", restful.DefaultHandler)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		co.serveWs(w, r)
	})

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := http.ListenAndServe(*addr, nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
}
