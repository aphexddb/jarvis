package client

import (
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// connection timeouts
const (
	MaxReconnectWait = time.Second * 30
)

// buffer sizes
const (
	SocketInChannelSize  = 100
	SocketOutChannelSize = 10
)

const maxConnectDelay = 5

type stateFn func(s *Socket) stateFn

// Socket represents a client that communicates with a Home Service
type Socket struct {
	deviceID        string
	client          Client
	url             string
	in              chan []byte
	out             chan []byte
	ws              *websocket.Conn
	closeSocket     chan bool
	lastActivity    time.Time
	connectTimeout  time.Duration
	activityTimeout time.Duration
	pingTimeout     time.Duration
	connectDelay    time.Duration
	timeoutTimer    *TimeoutTimer
}

// getConnectDelay returns a random connect delay to prevent thundering herds
func getConnectDelay() time.Duration {
	var rs = rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(rs)
	seconds := r1.Intn(maxConnectDelay)
	duration := time.Second * time.Duration(seconds+1)
	log.Println("Using random connect delay of ", duration)
	return duration
}

// NewSocket returns a new websocket
func NewSocket(u url.URL, cf Config, client Client) *Socket {
	s := &Socket{
		client:          client,
		url:             u.String(),
		connectTimeout:  cf.ConnectTimeout,
		activityTimeout: cf.ActivityTimeout,
		pingTimeout:     cf.PingTimeout,
		connectDelay:    0,
		out:             make(chan []byte, SocketOutChannelSize),
		closeSocket:     make(chan bool),
		timeoutTimer:    NewTimeoutTimer(NoTimeout, 0),
	}
	return s
}

// SetTimeout sets a timeout on a websocket
func (s *Socket) SetTimeout(reason TimeoutReason, d time.Duration) {
	s.timeoutTimer.SetTimeout(reason, d)
}

// startState is the starting state function
func startState(s *Socket) stateFn {
	s.connectDelay = getConnectDelay()

	// handle delayed re-connects
	if s.connectDelay > 0 {
		if s.connectDelay > MaxReconnectWait {
			s.connectDelay = MaxReconnectWait
		}
		time.Sleep(s.connectDelay)
	}
	// Set connection timeout on dialer
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = s.connectTimeout
	s.SetTimeout(ConnectTimeout, s.connectTimeout)
	ws, _, err := dialer.Dial(s.url, nil)
	if err != nil {
		log.Println("Error connecting:", err)
		// increase delay & reconnect
		s.connectDelay += time.Second
		return startState
	}
	// websocket connected
	s.ws = ws
	// Start reader & writer
	s.makeReader()
	s.makeWriter()
	s.client.HandleConnected()
	return connectedState
}

// Close closes a websocket connection
func (s *Socket) Close() {
	s.closeSocket <- true
}

// updateActivity updates the timeout timer
func (s *Socket) updateActivity() {
	s.timeoutTimer.Reset()
}

// connectedState is connection successful state function
func connectedState(s *Socket) stateFn {
	for {
		// wait for event from reader or heartbeat
		select {
		case event := <-s.in:
			if event == nil {
				return reconnectState
			}
			s.updateActivity()
			if err := s.client.HandleMessage(event); err != nil {
				return s.errorState(err)
			}
		case tick := <-s.timeoutTimer.C:
			// process timeout timer ticks
			if s.timeoutTimer.tickExpired(tick) {
				return timeoutState
			}
			continue
		case <-s.closeSocket:
			return stopState
		}
	}
}

// reconnectState is the reconnect state function
func reconnectState(s *Socket) stateFn {
	s.reset()
	if s.client.HandleDisconnect() {
		return startState
	}
	return nil
}

// timeoutState is the timeout state function
func timeoutState(s *Socket) stateFn {
	switch s.timeoutTimer.Reason {
	case ActivityTimeout:
		s.sendPing()
	case ConnectTimeout:
		log.Println("Connect timeout.")
		return reconnectState
	case PingTimeout:
		log.Println("Ping timeout.")
		return reconnectState
	}
	return connectedState
}

// stopState is the stop state function
func stopState(s *Socket) stateFn {
	err := s.ws.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
	}

	// close socket
	if s.ws != nil {
		s.ws.Close()
		s.ws = nil
	}

	return nil
}

// errorState is the error state function
func (s *Socket) errorState(err error) stateFn {
	log.Println("Websocket error:", err)
	switch err := err.(type) {
	case DelayError:
		if err.Temporary() {
			delay := err.Delay()
			if delay > 0 {
				s.connectDelay += delay
			} else {
				s.connectDelay = 0
			}
			return reconnectState
		}
	}
	return stopState
}

// reset is the reset state function
func (s *Socket) reset() {
	if s.out != nil {
		close(s.out)
		// create a new out channel
		s.out = make(chan []byte, SocketOutChannelSize)
	}
	if s.ws != nil {
		s.ws.Close()
		s.ws = nil
	}
}

// sendPing sends a ping over the websocket
func (s *Socket) sendPing() {
	// set ping timeout
	s.SetTimeout(PingTimeout, s.pingTimeout)
	// send ping
	s.client.SendPing()
}

// run starts monitoring for incoming messages
func (s *Socket) run() {
	for state := startState; state != nil; {
		state = state(s)
	}
}
