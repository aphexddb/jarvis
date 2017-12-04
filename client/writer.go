package client

import (
	"log"

	"github.com/gorilla/websocket"
)

// SendMessage sends a message over a websocket
func (s *Socket) SendMessage(msg []byte) {
	s.out <- msg
}

// makeWriter creates a writer goroutine
func (s *Socket) makeWriter() {
	ws := s.ws
	out := s.out
	go func() {
		for {
			msg, ok := <-out
			if !ok {
				// stop writer
				return
			}
			err := ws.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("Writer error:", err)
				return
			}
		}
	}()
}
