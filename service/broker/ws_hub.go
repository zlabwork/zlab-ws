// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package broker

import (
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {

	// Registered clients.
	clients map[int64]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[int64]*Client),
	}
}

func (h *Hub) GetClientsNumber() int {
	return len(h.clients)
}

func (h *Hub) Run() {

	var enqueued, producerErrors int
	producer := getProducer(nil)
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client

		case client := <-h.unregister:
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
			}

		case message := <-h.broadcast:
			producer.Input() <- &sarama.ProducerMessage{Topic: "message", Key: nil, Value: sarama.ByteEncoder(message)}
			enqueued++

		case err := <-producer.Errors():
			log.Println("Failed to produce message", err)
			producerErrors++
		}
	}

}
