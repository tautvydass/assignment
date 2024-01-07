package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"assignment/server/config"
	"assignment/server/server"

	"go.uber.org/zap"
)

func main() {
	// Resolve config path from command line arguments.
	args := os.Args[1:]
	if len(args) == 0 {
		panic("missing (config file, cert file, key file) argument(s)")
	}
	if len(args) != 3 {
		panic("mismatching number of arguments, expected 3 (config file, cert file, key file)")
	}
	path := args[0]

	// Load config from given path.
	config, err := config.LoadConfig(path)
	if err != nil {
		panic(fmt.Sprintf("error loading config at path %q: %v", path, err))
	}

	// Load TLS certificate and key.
	cert, err := tls.LoadX509KeyPair(args[1], args[2])
	if err != nil {
		panic(fmt.Sprintf("error loading TLS certificate and key: %v", err))
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("error creating logger: %v", err))
	}
	defer logger.Sync()

	// Start the server.
	logger.Info("Starting server")
	server := server.New(server.Config{
		SubscriberPort:    config.SubscriberPort,
		PublisherPort:     config.PublisherPort,
		TLS:               tlsConfig,
		OpenStreamTimeout: config.OpenStreamTimeout,
	}, logger)
	if err := server.Start(); err != nil {
		panic(fmt.Sprintf("error starting server: %v", err))
	}
	logger.Info("Server started and running")

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
		logger.Info("Graceful shutdown complete")
		return
	case <-ctx.Done():
		logger.Info("Graceful shutdown timed out, shutting down")
		return
	}
}
