package client

import (
	"time"
)

// TimeoutReason defines why a websocket times out
type TimeoutReason int

// timeout reasons
const (
	NoTimeout TimeoutReason = iota
	ConnectTimeout
	ActivityTimeout
	PingTimeout
)

// TimeoutTimer represents a websocket connection timer
type TimeoutTimer struct {
	C        <-chan time.Time
	Reason   TimeoutReason
	ticker   *time.Ticker
	duration time.Duration
	start    time.Time
}

func (t *TimeoutTimer) tickExpired(tick time.Time) bool {
	// check if timeout is disabled.
	if t.Reason == NoTimeout {
		// ignore tick
		return false
	}
	// check if the timeout has expired
	expire := tick.Sub(t.start)
	if expire >= t.duration {
		return true
	}
	return false
}

// SetTimeout sets a timeout and resets the timer
func (t *TimeoutTimer) SetTimeout(reason TimeoutReason, d time.Duration) {
	t.Reason = reason
	t.duration = d
	t.Reset()
}

// Reset resets the timer
func (t *TimeoutTimer) Reset() {
	t.start = time.Now()
}

// Expired returns if the timer has expired
func (t *TimeoutTimer) Expired() bool {
	// process all buffered ticks from ticker.
	for {
		select {
		case tick := <-t.C:
			if t.tickExpired(tick) {
				return true
			}
		default:
			// no more ticks
			return false
		}
	}
}

// Stop stops the timer
func (t *TimeoutTimer) Stop() {
	t.ticker.Stop()
}

// NewTimeoutTimer returns a new timeout timer
func NewTimeoutTimer(reason TimeoutReason, d time.Duration) *TimeoutTimer {
	t := &TimeoutTimer{
		Reason:   reason,
		ticker:   time.NewTicker(time.Second),
		duration: d,
		start:    time.Now(),
	}
	t.C = t.ticker.C
	return t
}
