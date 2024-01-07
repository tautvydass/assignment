package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"assignment/lib/receiver"

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

	// Set up the receiver.
	connectionClosed := make(chan struct{})
	receiver := receiver.New(logger)
	if err := receiver.Start(port, connectionClosed); err != nil {
		panic(fmt.Sprintf("error starting receiver: %v", err))
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

	if err := receiver.Close(); err != nil {
		logger.Error("Error closing receiver", zap.Error(err))
	}
}
