package service

import (
	"app"
	pb "app/grpc/monitor"
	"app/restful"
	"app/service/broker"
	"app/service/cache"
	"app/service/control"
	"context"
	"log"
	"net/http"
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

func (br *Broker) monitor() {

	type info struct {
		Node    string `json:"node"`
		Time    int64  `json:"time"`
		Clients int    `json:"clients"`
	}

	ticker := time.NewTicker(8 * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C

		go func() {

			// TODO: get config from Yaml
			node := app.Yaml.Base.Node
			num := int32(br.hub.GetClientsNumber())

			// conn
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			conn, err := control.MonitorConn()
			if err != nil {
				return
			}
			defer conn.Close()
			defer cancel()

			// TODO: transfer data
			cli := pb.NewMonitorClient(conn)
			_, err = cli.Notice(ctx, &pb.BrokerData{Id: node, Number: num})
			if err != nil {
				log.Println(err)
				return
			}
			_, err = cli.Health(ctx, &pb.HealthData{Id: 1, Ip: app.Yaml.Base.Host, Role: "Broker"})
			if err != nil {
				log.Println(err)
				return
			}
		}()

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
	go br.monitor()

	// websocket
	log.Println("broker service at " + *addr)
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
