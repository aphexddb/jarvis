package client

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aphexddb/jarvis"
	"github.com/aphexddb/jarvis/hardware"
	"periph.io/x/periph/conn/gpio"
)

// Client defines a websocket client
type Client interface {
	// return true to reconnect
	HandleDisconnect() bool
	HandleConnected()

	HandleMessage([]byte) error

	SendMessage([]byte)

	SendPing()

	Close()
}

// ChannelClient defines a channel connection
type ChannelClient interface {
	Client

	SendEvent(event Event)

	SendSubscribe(channel string)
	SendUnsubscribe(channel string)

	Subscribe(channel string) Channel
	Unsubscribe(channel string)
}

// ServiceClient represents a websocket client
type ServiceClient struct {
	deviceID string
	socket   *Socket
	hw       *hardware.RPI
}

// NewServiceClient returns a new ServiceClient
func NewServiceClient(deviceID string) *ServiceClient {
	return &ServiceClient{
		deviceID: deviceID,
		hw:       hardware.NewRPI(),
	}
}

// StartWithSocket starts a client
func (c *ServiceClient) StartWithSocket(socket *Socket) {
	c.socket = socket
	go c.socket.run()
}

// register performs a registration action
func (c *ServiceClient) register() {
	msg := jarvis.ClientMessage{
		Action: jarvis.ActionRegister,
		ID:     c.deviceID,
	}
	bytes, _ := json.Marshal(msg)
	c.SendMessage(bytes)
}

// HandleDisconnect handles a disconnect event
func (c *ServiceClient) HandleDisconnect() bool {
	log.Println("[Action] disconnected")
	return true
}

// HandleConnected handles a connection event
func (c *ServiceClient) HandleConnected() {
	log.Println("[Action] connected")
	c.register()
}

// HandleMessage handles a message event
func (c *ServiceClient) HandleMessage(msg []byte) error {
	log.Println("[Action] handle message: ", string(msg))
	return c.parseServerMessage(msg)
}

// SendMessage handles sending a message
func (c *ServiceClient) SendMessage(msg []byte) {
	log.Println("[Action] send message: ", string(msg))
	c.socket.SendMessage(msg)
}

// SendPing sends a ping
func (c *ServiceClient) SendPing() {
	log.Println("[Action] send ping")
	c.socket.sendPing()

}

// Close closes a connection to a server
func (c *ServiceClient) Close() {
	log.Println("[State] close")
	c.socket.closeSocket <- true
}

// parseServerMessage reads a message from the server
func (c *ServiceClient) parseServerMessage(msg []byte) error {
	sm := jarvis.ServerMessage{}
	err := json.Unmarshal(msg, &sm)
	if err != nil {
		return fmt.Errorf("recv: unknown message format: %s\n%s", err.Error(), string(msg))
	}

	// handle message types
	switch sm.Action {
	default:
		log.Println("recv:", sm.Message)
	}

	// toogle PIN
	log.Println("Testing, received a message and toggling Pin 33")
	if c.hw.GetState().GetGPIOPin33() {
		c.hw.SetGPIOPin33(gpio.Low)
	} else {
		c.hw.SetGPIOPin33(gpio.High)
	}

	return nil
}
