package server

import (
	"context"
	"crypto/tls"
	"time"

	"assignment/lib/connection"
	"assignment/lib/log"
	"assignment/server/server/controller"
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
	SubscriberPort     int
	PublisherPort      int
	TLS                *tls.Config
	OpenStreamTimeout  time.Duration
	SendMessageTimeout time.Duration
}

// Server is an interface for the broker server.
type Server interface {
	Start() error
	Shutdown() error
}

// New creates a new broker server.
func New(config Config) Server {
	return &server{
		config:          config,
		newListener:     listener.New,
		commsController: controller.NewCommsController(),
	}
}

type server struct {
	config             Config
	started            bool
	publisherListener  listener.Listener
	subscriberListener listener.Listener
	commsController    controller.CommsController

	// listener constructor delegate used for mocks
	newListener func(cb listener.NewConnectionCallback) listener.Listener
}

func (s *server) Start() error {
	if s.started {
		return ErrAlreadyStarted
	}

	s.publisherListener = s.newListener(s.addPublisher)
	if err := s.publisherListener.Start(s.config.PublisherPort, s.config.TLS); err != nil {
		return errors.Wrap(err, "start publisher listener")
	}
	log.Tracef("Started publisher listener on port %d", s.config.PublisherPort)

	s.subscriberListener = s.newListener(s.addSubscriber)
	if err := s.subscriberListener.Start(s.config.SubscriberPort, s.config.TLS); err != nil {
		return errors.Wrap(err, "start subscriber listener")
	}
	log.Tracef("Started subscriber listener on port %d", s.config.SubscriberPort)

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

	if err := s.commsController.Close(); err != nil {
		return errors.Wrap(err, "close comms controller")
	}

	s.started = false
	return nil
}

func (s *server) addPublisher(conn connection.Connection) {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.OpenStreamTimeout)
	defer cancel()

	log.Trace("Publisher connected, opening read write stream")
	readWriteStream, err := conn.OpenReadWriteStream(ctx, s.commsController.MessageReceiver())
	if err != nil {
		log.Errorf("Error opening publisher stream: %s", err.Error())
		return
	}
	readWriteStream.SetSendMessageTimeout(s.config.SendMessageTimeout)

	s.commsController.AddPublisher(readWriteStream)
}

func (s *server) addSubscriber(conn connection.Connection) {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.OpenStreamTimeout)
	defer cancel()

	log.Trace("Subscriber connected, opening write stream")
	writeStream, err := conn.OpenWriteStream(ctx)
	if err != nil {
		log.Errorf("Error opening subscriber stream: %s", err.Error())
		return
	}
	writeStream.SetSendMessageTimeout(s.config.SendMessageTimeout)

	s.commsController.AddSubscriber(writeStream)
}
