package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aphexddb/jarvis/hardware"
	"periph.io/x/periph/conn/gpio"
)

var hw *hardware.RPI

const tickTimeMs = 1000

func main() {
	log.SetFlags(1)

	hw := hardware.NewRPI()

	log.Println("Ticking every", tickTimeMs, "ms")

	ticker := time.NewTicker(time.Millisecond * time.Duration(tickTimeMs))
	go func() {
		for range ticker.C {
			log.Println("PIN 33 state is", hw.GetState().GetGPIOPin33())
			if hw.GetState().GetGPIOPin33() {
				hw.SetGPIOPin33(gpio.Low)
			} else {
				hw.SetGPIOPin33(gpio.High)
			}
		}
	}()

	time.Sleep(time.Millisecond * 5000)
	ticker.Stop()
	fmt.Println("Ticker stopped")

}
