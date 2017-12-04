package client

import (
	"net"
	"time"
)

// DelayError defines a delay error state for a websocket
type DelayError interface {
	net.Error
	Delay() time.Duration // Delay reconnect?
}

// Error represents an error
type Error struct {
	reason  string
	timeout bool
	temp    bool
	delay   time.Duration
}

// Error returns an error reason
func (e *Error) Error() string {
	return e.reason
}

// Tiemout returns if there was a timeout
func (e *Error) Tiemout() bool {
	return e.timeout
}

// Temporary returns if the error was temporary
func (e *Error) Temporary() bool {
	return e.temp
}

// Delay returns the length of the delay
func (e *Error) Delay() time.Duration {
	return e.delay
}

// NewError returns a new error
func NewError(reason string, timeout bool, temp bool, delay time.Duration) *Error {
	return &Error{
		reason:  reason,
		timeout: timeout,
		temp:    temp,
		delay:   delay,
	}
}

// client errors
var (
	ErrReconnect      = NewError("Reconnect (no delay)", false, true, 0)
	ErrDelayReconnect = NewError("Reconnect (with delay)", false, true, time.Second)
	ErrClosed         = NewError("Closed", false, false, 0)
)
