package main

import (
	"bytes"
	"log"
	"net/http"
	
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	// register client
	clients[ws] = true	
}

func writer(t string, b []byte) {
	var buffer bytes.Buffer
	buffer.Write([]byte(t))
	buffer.Write([]byte("||"))
	buffer.Write(b)
	broadcast <- buffer.Bytes()
}

func echo() {
	for {
		b := <-broadcast
		// send to every client that is currently connected
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				log.Printf("Websocket error: %s", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
