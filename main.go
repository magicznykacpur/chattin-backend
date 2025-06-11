package main

import (
	"log"
	"net/http"

	"github.com/magicznykacpur/chattin/pkg/websocket"
)

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	log.Println("Websocket endpoint hit", r.Host)
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		log.Println(err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
	
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

func main() {
	log.Println("Setting up websocket...")
	setupRoutes()
	
	log.Println("Starting the server on port :42069...")

	err := http.ListenAndServe(":42069", nil)
	if err != nil {
		log.Fatal(err)
	}
}
