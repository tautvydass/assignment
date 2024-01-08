package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"assignment/client/subscriber/client"
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

	// Set up the subscriber client.
	connectionClosed := make(chan struct{})
	client := client.New()
	if err := client.Start(port, connectionClosed); err != nil {
		panic(fmt.Sprintf("error starting subscriber client: %v", err))
	}

	// Wait for SIGINT or SIGTERM.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	signal.Notify(shutdown, syscall.SIGTERM)

	select {
	case <-shutdown:
		log.Trace("Shutting down")
	case <-connectionClosed:
		log.Trace("Connection closed by the server, shutting down")
	}

	if err := client.Close(); err != nil {
		log.Errorf("Error closing subscriber client: %s", err.Error())
	}
}
