package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/aphexddb/jarvis/client"
)

var ssl = flag.Bool("ssl", false, "Connect using HTTPS")
var addr = flag.String("addr", "jarvis.yourdomain.com", "Service address")
var path = flag.String("path", "/v1/ws", "Service path")
var deviceID = flag.String("device", "", "Device id to register (uuid)")

func main() {
	flag.Parse()
	log.SetFlags(0)

	// make sure device id is set
	if len(*deviceID) != 36 {
		log.Fatal("Device ID missing or not, use the -device flag to set the device ID.\nYou can create a unique uuid here: https://www.uuidgenerator.net/.\n")
		os.Exit(1)
	}
	log.Println("Using device ID", *deviceID)

	// create websocket URL
	var u url.URL
	if *ssl {
		u = url.URL{Scheme: "wss", Host: *addr, Path: *path}
	} else {
		u = url.URL{Scheme: "ws", Host: *addr, Path: *path}
	}
	log.Println("Connecting to server", u.String())

	c := client.NewServiceClient(*deviceID)
	s := client.NewSocket(u, client.DefaultConfig, c)

	// listen for interrupts
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// main event loop
	c.StartWithSocket(s)

	func() {
		<-interrupt
		log.Println("interrupt received, attempting graceful close")
		// To cleanly close a connection, a client should send a close
		// frame and wait for the server to close the connection.
		go s.Close()
		time.Sleep(1 * time.Second)
		os.Exit(1)

	}()
}
