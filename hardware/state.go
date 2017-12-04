package hardware

// State represents a hardware state
type State struct {
	GPIOPin33 bool
}

// NewState returns a new hardware State
func NewState() State {
	return State{
		GPIOPin33: false,
	}
}

// GetGPIOPin33 returns the state of GPIO Pin 33
func (s State) GetGPIOPin33() bool {
	return s.GPIOPin33
}
