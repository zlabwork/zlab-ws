// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
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

			fmt.Println(message)

			msgType := message[1]
			msgLength := binary.BigEndian.Uint16(message[2:4])
			msgId := hex.EncodeToString(message[4:20])
			msgSender := int64(binary.BigEndian.Uint64(message[12:20]))
			msgReceiver := int64(binary.BigEndian.Uint64(message[20:28]))
			msgBody := string(message[28:])

			fmt.Println(msgType, msgLength, msgId)
			fmt.Println(msgSender, msgReceiver, msgBody)

			// send to user
			cli, ok := h.clients[msgReceiver]
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
