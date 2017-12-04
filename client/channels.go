package client

import (
	"sync"
)

// Channels represents a connection to a websocket
type Channels struct {
	sync.RWMutex
	client    ChannelClient
	channels  map[string]Channel
	global    Channel
	connected bool
}

// HandleEvent handles an event
func (c *Channels) HandleEvent(event Event) {
	// send event to global channel
	if c.global != nil {
		c.global.HandleEvent(event)
	}
	channelName := event.GetChannel()
	if channelName == "" {
		// global only event.
		return
	}
	// send event to subscribed channel
	ch := c.Find(channelName)
	if ch != nil {
		ch.HandleEvent(event)
	}
}

// ConnectedState updates the connected state of a websocket
func (c *Channels) ConnectedState(connected bool) {
	c.RLock()
	defer c.RUnlock()
	// cache connected state
	c.connected = connected
	// notify all channels of the connection state
	for _, ch := range c.channels {
		ch.UpdateClientState(connected)
	}
}

// SubscriptionSucceded notifies when a channel is opened sucessfully
func (c *Channels) SubscriptionSucceded(channel string, succeded bool) {
	c.RLock()
	defer c.RUnlock()
	ch := c.Find(channel)
	if ch != nil {
		ch.SetActive(true)
	}
}

// Find finds a channel by name
func (c *Channels) Find(channel string) Channel {
	// empty channel name is for receiving events from all subscribed channels.
	if channel == "" {
		return c.global
	}
	// mutex is for the 'channels' map
	c.RLock()
	defer c.RUnlock()
	return c.channels[channel]
}

// Add adds a channel by name
func (c *Channels) Add(channel string, ch Channel) {
	// empty channel name is for receiving events from all subscribed channels.
	if channel == "" {
		c.global = ch
		return
	}
	// mutex is for the 'channels' map
	c.Lock()
	defer c.Unlock()
	c.channels[channel] = ch
	if c.connected {
		ch.Subscribe()
	}
}

// Remove removes a channel by name
func (c *Channels) Remove(channel string) {
	// empty channel name is for receiving events from all subscribed channels.
	if channel == "" {
		c.global = nil
	}
	// mutex is for the 'channels' map
	c.Lock()
	defer c.Unlock()
	ch := c.channels[channel]
	if ch != nil && c.connected {
		ch.Unsubscribe()
	}
	delete(c.channels, channel)
}

// Bind binds an event handler to an event type
func (c *Channels) Bind(event string, h Handler) {
	c.global.Bind(event, h)
}

// Unbind unbinds an event handler to an event type
func (c *Channels) Unbind(event string, h Handler) {
	c.global.Unbind(event, h)
}

// NewChannels creates a new Channels struct
func NewChannels(client ChannelClient) *Channels {
	return &Channels{
		client:   client,
		channels: make(map[string]Channel),
	}
}
