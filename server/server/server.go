package server

import (
	"context"
	"crypto/tls"
	"time"

	"assignment/lib/connection"
	"assignment/lib/entity"
	"assignment/server/server/listener"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	// ErrNotStarted is returned when attempting to start a
	// server that's already started.
	ErrAlreadyStarted = errors.New("server already started")
)

// Config contains configuration for the broker server.
type Config struct {
	SubscriberPort    int
	PublisherPort     int
	TLS               *tls.Config
	OpenStreamTimeout time.Duration
}

// Server is an interface for the broker server.
type Server interface {
	Start() error
	Shutdown() error
}

// New creates a new broker server.
func New(config Config, logger *zap.Logger) Server {
	return &server{
		config:      config,
		logger:      logger,
		newListener: listener.New,
	}
}

type server struct {
	config  Config
	started bool
	logger  *zap.Logger

	publisherListener  listener.Listener
	subscriberListener listener.Listener

	// listener constructor delegate used for mocks
	newListener func(
		cb listener.NewConnectionCallback,
		logger *zap.Logger,
	) listener.Listener
}

func (s *server) Start() error {
	if s.started {
		return ErrAlreadyStarted
	}

	s.publisherListener = s.newListener(s.addPublisher, s.logger)
	if err := s.publisherListener.Start(s.config.PublisherPort, s.config.TLS); err != nil {
		return errors.Wrap(err, "start publisher listener")
	}
	s.logger.Info("Started publisher listener", zap.Int("port", s.config.PublisherPort))

	s.subscriberListener = s.newListener(s.addSubscriber, s.logger)
	if err := s.subscriberListener.Start(s.config.SubscriberPort, s.config.TLS); err != nil {
		return errors.Wrap(err, "start subscriber listener")
	}
	s.logger.Info("Started subscriber listener", zap.Int("port", s.config.SubscriberPort))

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
	s.logger.Info("Publisher connected", zap.Any("conn", conn))
}

func (s *server) addSubscriber(conn connection.Connection) {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.OpenStreamTimeout)
	defer cancel()

	s.logger.Info("Subscriber connected, opening write stream")
	writeStream, err := conn.OpenWriteStream(ctx)
	if err != nil {
		s.logger.Error("Error opening write stream", zap.Error(err))
		return
	}

	message := entity.Message{
		Text: "Hello from server!",
	}
	s.logger.Info("Sending message to subscriber", zap.String("message", message.Text))
	if err := writeStream.SendMessage(message); err != nil {
		s.logger.Error("Error sending message to subscriber", zap.Error(err))
	}

	// TODO: add the subscriber to communication controller

	// TODO: wait for the subscriber to receive the message and close the stream for now
	time.Sleep(time.Second)
	if err = writeStream.CloseStream(); err != nil {
		s.logger.Error("Error closing subscriber stream", zap.Error(err))
	}
}
