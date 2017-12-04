package jarvis

import "log"

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// NewHub returns a new Hub
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	log.Println("Starting Hub")
	for {
		select {

		case client := <-h.register:
			log.Println("Hub registered new client:", client.conn.RemoteAddr().String())
			h.clients[client] = true

		case client := <-h.unregister:
			log.Println("Hub unregistered client:", client.conn.RemoteAddr().String(), "("+client.id+")")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			ignored := 0
			log.Println("Hub broadcasting to registered clients:", string(message))

			for client := range h.clients {
				// only broadcast to clients that have registered with an ID
				if client.id != "" {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				} else {
					ignored = ignored + 1
				}
			}

			if ignored > 0 {
				log.Println("Hub broadcast message ignored", ignored, "non-registered clients")
			}
		}
	}
}
