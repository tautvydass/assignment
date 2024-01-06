package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"assignment/server/config"
	"assignment/server/server"
)

func main() {
	// Resolve config path from command line arguments.
	args := os.Args[1:]
	if len(args) == 0 {
		panic("missing config path argument")
	}
	if len(args) > 1 {
		panic("too many arguments")
	}
	path := args[0]

	// Load config from given path.
	config, err := config.LoadConfig(path)
	if err != nil {
		panic(fmt.Sprintf("error loading config at path %q: %v", path, err))
	}

	// Start the server.
	fmt.Println("Starting server")
	server := server.New(server.Config{
		SubscriberPort: config.SubscriberPort,
		PublisherPort:  config.PublisherPort,
	})
	if err := server.Start(); err != nil {
		panic(fmt.Sprintf("error starting server: %v", err))
	}
	fmt.Println("Server started and running")

	// Set up graceful shutdown.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	signal.Notify(shutdown, syscall.SIGTERM)

	// Wait for shutdown signal and shutdown the server.
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), config.GracefulShutdownTimeout)
	defer cancel()
	closed := make(chan struct{})
	go func() {
		if err := server.Shutdown(); err != nil {
			panic(fmt.Sprintf("error shutting down server: %v", err))
		}
		closed <- struct{}{}
	}()

	select {
	case <-closed:
		fmt.Println("Graceful shutdown complete")
		return
	case <-ctx.Done():
		fmt.Println("Graceful shutdown timed out, shutting down")
		return
	}
}
