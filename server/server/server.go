package server

import (
	"fmt"

	"assignment/lib/connection"
	"assignment/server/server/listener"

	"github.com/pkg/errors"
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
		config:      config,
		newListener: listener.New,
	}
}

type server struct {
	config  Config
	started bool

	publisherListener  listener.Listener
	subscriberListener listener.Listener

	// listener constructor delegate used for mocks
	newListener func(cb listener.NewConnectionCallback) listener.Listener
}

func (s *server) Start() error {
	if s.started {
		return ErrAlreadyStarted
	}

	s.publisherListener = s.newListener(s.addPublisher)
	if err := s.publisherListener.Start(s.config.PublisherPort); err != nil {
		return errors.Wrap(err, "start publisher listener")
	}
	fmt.Printf("Started publisher listener on port %d\n", s.config.PublisherPort)

	s.subscriberListener = s.newListener(s.addSubscriber)
	if err := s.subscriberListener.Start(s.config.SubscriberPort); err != nil {
		return errors.Wrap(err, "start subscriber listener")
	}
	fmt.Printf("Started subscriber listener on port %d\n", s.config.SubscriberPort)

	s.started = true
	return nil
}

func (s *server) Shutdown() error {
	if !s.started {
		return nil
	}

	if err := s.publisherListener.Shutdown(); err != nil {
		return errors.Wrap(err, "shutdown publisher listener")
	}

	if err := s.subscriberListener.Shutdown(); err != nil {
		return errors.Wrap(err, "shutdown subscriber listener")
	}

	s.started = false
	return nil
}

func (s *server) addPublisher(conn connection.Connection) {
	// TODO: create a bidirectional stream
	// TODO: save the publisher in memory
	fmt.Printf("Publisher connected: %v\n", conn)
}

func (s *server) addSubscriber(conn connection.Connection) {
	// TODO: create a simple stream
	// TODO: save the subscriber in memory
	fmt.Printf("Pubscriber connected: %v\n", conn)
}
