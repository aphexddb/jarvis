// +build darwin

package hardware

import (
	"log"

	"periph.io/x/periph/conn/gpio"
)

// init initilizes an rPi
func (pi *RPI) init() {
	log.Println("[MOCK] Raspberry PI initalizing")
}

// SetGPIOPin33 sets GPIO pin 33
func (pi *RPI) SetGPIOPin33(level gpio.Level) (err error) {
	log.Println("[MOCK] Setting GPIO pin 33", level)
	pi.state.GPIOPin33 = levelToBool(level)
	return err
}
