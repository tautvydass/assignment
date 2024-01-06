package server

import (
	"crypto/tls"

	"assignment/lib/connection"
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
	SubscriberPort int
	PublisherPort  int
}

// Server is an interface for the broker server.
type Server interface {
	Start() error
	Shutdown() error
}

// New creates a new broker server.
func New(
	config Config,
	logger *zap.Logger,
	tlsConfig *tls.Config,
) Server {
	return &server{
		config:      config,
		logger:      logger,
		tlsConfig:   tlsConfig,
		newListener: listener.New,
	}
}

type server struct {
	config    Config
	started   bool
	logger    *zap.Logger
	tlsConfig *tls.Config

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
	if err := s.publisherListener.Start(s.config.PublisherPort, s.tlsConfig); err != nil {
		return errors.Wrap(err, "start publisher listener")
	}
	s.logger.Info("Started publisher listener", zap.Int("port", s.config.PublisherPort))

	s.subscriberListener = s.newListener(s.addSubscriber, s.logger)
	if err := s.subscriberListener.Start(s.config.SubscriberPort, s.tlsConfig); err != nil {
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
	// TODO: create a simple stream
	// TODO: save the subscriber in memory
	s.logger.Info("Subscriber connected", zap.Any("conn", conn))
}
