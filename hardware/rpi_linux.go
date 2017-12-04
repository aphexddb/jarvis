// +build linux

package hardware

import (
	"log"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/rpi"
)

// init initilizes an rPi
func (pi *RPI) init() {
	log.Println("Raspberry PI initalizing")

	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
}

// SetGPIOPin33 sets GPIO pin 33
func (pi *RPI) SetGPIOPin33(level gpio.Level) (err error) {
	log.Println("Setting GPIO pin 33", level)

	if level {
		err = rpi.P1_33.Out(gpio.High)
	} else {
		err = rpi.P1_33.Out(gpio.Low)
	}
	pi.state.GPIOPin33 = levelToBool(level)

	return err
}
