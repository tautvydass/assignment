package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"assignment/lib/connection"
	"assignment/lib/entity"

	"github.com/pkg/errors"
	"github.com/quic-go/quic-go"
	"go.uber.org/zap"
)

const (
	// DefaultTimeout is the default timeout for establishing
	// a connection to the server.
	DefaultTimeout = time.Second * 30
	// DefaultIdleTimeout is the default idle timeout for the
	// connection.
	DefaultIdleTimeout = time.Hour
)

// Client is an interface for publishing messages to the server. it
// will also log all messages received from the server.
type Client interface {
	// Start establishes a connection with the server and starts
	// listening to messages. Given channel is closed when connection.
	Start(port int, connectionClosed chan struct{}) error
	// Publish publishes a message to the server.
	Publish(message string) error
	// Close closes the connection with the server.
	Close() error
}

type client struct {
	logger *zap.Logger
	stream connection.ReadWriteStream
}

// New constructs a new publisher client.
func New(logger *zap.Logger) Client {
	return &client{
		logger: logger,
	}
}

func (c *client) Start(port int, connectionClosed chan struct{}) error {
	var err error
	c.stream, err = c.setupReadWriteStream(port, connectionClosed)
	if err != nil {
		return errors.Wrap(err, "setup read write stream")
	}

	c.logger.Info("Started listening to messages", zap.Int("port", port))
	return nil
}

func (c *client) setupReadWriteStream(
	port int, connectionClosed chan struct{},
) (connection.ReadWriteStream, error) {
	c.logger.Info("Setting up UDP connection")
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: 0})
	if err != nil {
		return nil, errors.Wrap(err, "listen udp")
	}
	transport := &quic.Transport{Conn: udpConn}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	address := fmt.Sprintf("localhost:%d", port)
	c.logger.Info("Dialing the server", zap.String("address", address))
	conn, err := transport.Dial(
		ctx, &net.UDPAddr{Port: port}, &tls.Config{InsecureSkipVerify: true}, &quic.Config{
			MaxIdleTimeout: DefaultIdleTimeout,
		})
	if err != nil {
		return nil, errors.Wrapf(err, "dial %q", address)
	}

	c.logger.Info("Accepting stream")
	return connection.New(conn).
		AcceptReadWriteStream(ctx, c.handleMessage, connectionClosed)
}

func (c *client) Publish(message string) error {
	if err := c.stream.SendMessage(entity.Message{
		Text: message,
	}); err != nil {
		return errors.Wrap(err, "send message")
	}
	return nil
}

func (c *client) handleMessage(message entity.Message) {
	c.logger.Info("Received message", zap.String("message", message.Text))
}

func (c *client) Close() error {
	if c.stream == nil {
		return nil
	}
	return errors.Wrap(
		c.stream.CloseStream(),
		"close stream",
	)
}
