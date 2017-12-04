package client

import (
	"time"
)

// connection timeout defaults
const (
	connectTimeout  = time.Second * 30
	activityTimeout = time.Second * 60
	pingTimeout     = time.Second * 30
)

// Config represents client configuration
type Config struct {
	ConnectTimeout  time.Duration
	ActivityTimeout time.Duration
	PingTimeout     time.Duration
}

// DefaultConfig contains default config values
var DefaultConfig = Config{
	ConnectTimeout:  connectTimeout,
	ActivityTimeout: activityTimeout,
	PingTimeout:     pingTimeout,
}
