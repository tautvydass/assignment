package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"assignment/client/subscriber/receiver"
	"assignment/lib/log"
)

func main() {
	// Resolve port from command line arguments.
	args := os.Args[1:]
	if len(args) == 0 {
		panic("missing port argument")
	}
	if len(args) > 1 {
		panic("mismatching number of arguments, expected 1 (port)")
	}

	port, err := strconv.Atoi(args[0])
	if err != nil {
		panic(fmt.Sprintf("error parsing port: %v", err))
	}

	// Set up the receiver.
	connectionClosed := make(chan struct{})
	receiver := receiver.New()
	if err := receiver.Start(port, connectionClosed); err != nil {
		panic(fmt.Sprintf("error starting receiver: %v", err))
	}

	// Wait for SIGINT or SIGTERM.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	signal.Notify(shutdown, syscall.SIGTERM)

	select {
	case <-shutdown:
		log.Trace("Shutting down")
	case <-connectionClosed:
		log.Info("Connection closed by the server, shutting down")
	}

	if err := receiver.Close(); err != nil {
		log.Errorf("Error closing receiver: %s", err.Error())
	}
}
