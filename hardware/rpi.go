package hardware

import "periph.io/x/periph/conn/gpio"

// RPIO defines a raspberry Pi
type RPIO interface {
	Init()
	TogglePin33(bool) error
}

// RPI represents an instance of a Raspberry Pi
type RPI struct {
	state State
}

// NewRPI creates a new Raspberry Pi hardware interface
func NewRPI() *RPI {
	rpi := RPI{
		state: NewState(),
	}
	rpi.init()
	return &rpi
}

// GetState returns current state
func (pi *RPI) GetState() State {
	return pi.state
}

// levelToBool returns true for High and false for Low
func levelToBool(level gpio.Level) bool {
	if level == gpio.High {
		return true
	}
	return false
}
