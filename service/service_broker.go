package service

import (
	"app"
	"app/restful"
	"app/service/broker"
	"app/service/cache"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Broker struct {
	cache app.CacheToken
	repo  app.RepoMessage
	hub   *broker.Hub
}

func NewBrokerService() (*Broker, error) {

	cs, err := cache.NewTokenRepository()
	if err != nil {
		return nil, err
	}
	hub := broker.NewHub()

	return &Broker{
		cache: cs,
		hub:   hub,
	}, nil
}

func (br *Broker) information() {

	type info struct {
		Node    string `json:"node"`
		Time    int64  `json:"time"`
		Clients int    `json:"clients"`
	}

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C

		bs, err := json.Marshal(&info{Node: os.Getenv("APP_NODE"), Time: time.Now().UTC().Unix(), Clients: br.hub.GetClientsNumber()})
		if err != nil {
			return
		}
		log.Println("INFO:", string(bs))
	}
}

func (br *Broker) serveWs(w http.ResponseWriter, r *http.Request) {
	broker.ServeWs(br.hub, br.cache, w, r)
}

func (br *Broker) Run(addr *string) {

	http.HandleFunc("/", restful.DefaultHandler)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		br.serveWs(w, r)
	})

	// Run our server in a goroutine so that it doesn't block.
	go br.hub.Run()
	go br.information()

	// websocket
	go func() {
		if err := http.ListenAndServe(*addr, nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	// grpc server
	go func() {
		broker.StartRPC()
	}()
}
