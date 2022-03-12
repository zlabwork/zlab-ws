// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"encoding/json"
	"log"
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

func (h *Hub) Run() {

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

			// send to user
			type who struct {
				From int64
				To   int64
			}
			var w who
			if err := json.Unmarshal(message[2:], &w); err != nil {
				log.Println(err.Error())
				continue
			}
			cli, ok := h.clients[w.To]
			if !ok {
				// TODO :: 存储到数据库
				log.Println("the receiver user is not online")
				continue
			}
			select {
			case cli.send <- message:
			default:
				close(cli.send)
				delete(h.clients, cli.id)
			}
			// TODO :: send to group or channel

		}
	}

}
