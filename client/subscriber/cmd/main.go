package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"assignment/client/subscriber/receiver"

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
	// TODO: handle the connection closed callback.
	receiver := receiver.New(logger)
	if err := receiver.Start(port); err != nil {
		panic(fmt.Sprintf("error starting receiver: %v", err))
	}

	// Wait for SIGINT or SIGTERM.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	signal.Notify(shutdown, syscall.SIGTERM)

	<-shutdown
	logger.Info("Shutting down")
	receiver.Close()
}
