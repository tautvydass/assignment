package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"assignment/client/publisher/client"

	"go.uber.org/zap"
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

	// Set up logger.
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("error creating logger: %v", err))
	}
	defer logger.Sync()

	// Set up the publisher client.
	connectionClosed := make(chan struct{})
	client := client.New(logger)
	if err := client.Start(port, connectionClosed); err != nil {
		panic(fmt.Sprintf("error starting publisher client: %v", err))
	}

	// TODO: implement a proper way to publish messages.
	time.Sleep(time.Second)
	if err := client.Publish("Hello everyone! I'm the new publisher!"); err != nil {
		panic(fmt.Sprintf("error publishing message: %v", err))
	}

	// Wait for SIGINT or SIGTERM.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	signal.Notify(shutdown, syscall.SIGTERM)

	select {
	case <-shutdown:
		logger.Info("Shutting down")
	case <-connectionClosed:
		logger.Info("Connection closed by the server, shutting down")
	}

	if err := client.Close(); err != nil {
		logger.Error("Error closing publisher client", zap.Error(err))
	}
}
