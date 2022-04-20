// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package broker

import (
	"app"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 20 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 2048

	// Message head size 4 bytes (32 Bit)
	headSize = 4

	// Message body part size 24 bytes (64 * 3 Bit)
	bodyHeadSize = 24
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	cache app.CacheFace

	repo app.RepoFace

	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// user id
	id int64

	// secret key TODO: delete unused
	key []byte

	// AES cipher.Block
	block cipher.Block
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		log.Println(fmt.Sprintf("off line, client id: %d", c.id))
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		// 1. Read
		_, data, err := c.conn.ReadMessage()
		if err != nil || len(data) < headSize {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// 2. Message Decryption
		msg, err := decrypt(c.block, data[headSize:])
		if err != nil {
			log.Println("error: decrypt the message")
			break
		}
		copy(data[headSize:], msg)

		// 3. switch
		switch data[1] {
		case app.TypeAuth:
			if !c.authorize(data[:len(msg)+headSize]) {
				log.Println(fmt.Errorf("authorization failed"))
				break
			}
			c.hub.register <- c
			c.sendCachedData()

		default:
			// TODO: 处理粘包问题
			// len := binary.BigEndian.Uint16(data[2:headSize])
			c.hub.broadcast <- data[:len(msg)+headSize]
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}

			w.Write(*c.encrypt(message))

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(*c.encrypt(<-c.send))
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// check authorization
func (c *Client) authorize(msg []byte) bool {

	var au app.MsgAuth
	err := json.Unmarshal(msg[headSize+bodyHeadSize:], &au)
	if err != nil {
		log.Println(err)
		return false
	}
	id, err := strconv.ParseInt(au.Sender, 10, 64)
	if err != nil {
		return false
	}

	// TODO: check token
	log.Println("SUCCESS:", au)
	c.id = id
	c.hub.register <- c
	return true
}

// send message which cached in database
func (c *Client) sendCachedData() {
	data, err := c.repo.GetTodo(c.id)
	if err != nil {
		return
	}

	for _, item := range data {
		rev := make([]byte, 8)
		binary.BigEndian.PutUint64(rev, uint64(item.Receiver))

		mid, err := hex.DecodeString(item.Mid)
		if err != nil {
			continue
		}

		b := make([]byte, 28+len(item.Data))
		copy(b[1:2], []byte{item.Type})
		copy(b[4:20], mid)
		copy(b[20:28], rev)
		copy(b[28:], item.Data)

		c.send <- b
	}
	// Delete message cached
	c.repo.DeleteTodo(c.id)
}

// Message Encryption
func (c *Client) encrypt(message []byte) *[]byte {
	ciphertext, err := encrypt(c.block, message[headSize:])
	if err != nil {
		return nil
	}
	l := len(ciphertext) + headSize
	lb := make([]byte, 16)
	binary.BigEndian.PutUint16(lb, uint16(l))
	data := make([]byte, l)
	copy(data[0:2], message[:2])
	copy(data[2:headSize], lb)
	copy(data[headSize:], ciphertext)
	return &data
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, cache app.CacheFace, repo app.RepoFace, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{cache: cache, repo: repo, hub: hub, conn: conn, send: make(chan []byte, 256)}

	// TODO: dev test, need to Delete
	b, _ := hex.DecodeString("ffffffffffffffffffffffffffffffff")
	client.block, _ = aes.NewCipher(b)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
