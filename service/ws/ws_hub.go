// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"app"
	"encoding/binary"
	"encoding/hex"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Database repository
	repo app.RepoFace

	// Registered clients.
	clients map[int64]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub(repo app.RepoFace) *Hub {
	return &Hub{
		repo:       repo,
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
			now := time.Now()

			msgType := message[1]
			msgId := hex.EncodeToString(message[4:20])
			msgSender := int64(binary.BigEndian.Uint64(message[12:20]))
			msgReceiver := int64(binary.BigEndian.Uint64(message[20:28]))
			// msgLength := binary.BigEndian.Uint16(message[2:4])
			// msgBody := string(message[28:])

			// TODO :: send to group or channel

			// Save to database
			go func() {
				if h.repo.CreateLogs(msgType, msgId, msgSender, msgReceiver, message[28:], now) != nil {
					return
				}
			}()

			// send to user
			cli, ok := h.clients[msgReceiver]
			if !ok {
				go func() {
					if h.repo.CreateTodo(msgType, msgId, msgSender, msgReceiver, message[28:], now) != nil {
						return
					}
				}()
				continue
			}
			select {
			case cli.send <- message:
			default:
				close(cli.send)
				delete(h.clients, cli.id)
			}
		}
	}

}
