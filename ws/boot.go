// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"flag"
	"log"
	"net/http"
	"os"
)

var addr = flag.String("addr", ":8080", "http service address")

func Run() {

	if len(os.Getenv("APP_PORT")) > 0 {
		*addr = ":" + os.Getenv("APP_PORT")
	}
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/dev", devHandler)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
