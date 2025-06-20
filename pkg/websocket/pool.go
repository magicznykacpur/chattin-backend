package websocket

import "log"

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {

		case client := <-pool.Register:
			pool.Clients[client] = true
			log.Println("Size of connection pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				log.Println(client)
				client.Conn.WriteJSON(Message{Type: 1, Body: "New user joined..."})
			}

		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			log.Println("Size of connection pool: ", len(pool.Clients))

			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User disconnected..."})
			}

		case message := <-pool.Broadcast:
			log.Println("Sending message to all clients in the pool...")
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}
