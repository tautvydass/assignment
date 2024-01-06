package server

import (
	"errors"
)

var (
	// ErrNotStarted is returned when attempting to start a
	// server that's already started.
	ErrAlreadyStarted = errors.New("server already started")
)

// Config contains configuration for the broker server.
type Config struct {
	SubscriberPort int
	PublisherPort  int
}

// Server is an interface for the broker server.
type Server interface {
	Start() error
	Shutdown() error
}

// New creates a new broker server.
func New(config Config) Server {
	return &server{
		config: config,
	}
}

type server struct {
	close   chan struct{}
	config  Config
	started bool
}

func (s *server) Start() error {
	if s.started {
		return ErrAlreadyStarted
	}

	// TODO: start the server

	s.close = make(chan struct{})
	s.started = true
	go s.run()

	return nil
}

func (s *server) run() {
	for {
		select {
		case <-s.close:
			return
		default:
			// TODO: accept connections
		}
	}
}

func (s *server) Shutdown() error {
	s.close <- struct{}{}
	// TODO: close connections and the server
	s.started = false
	return nil
}
