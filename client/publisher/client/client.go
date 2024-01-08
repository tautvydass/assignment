package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"assignment/lib/connection"
	"assignment/lib/entity"
	"assignment/lib/log"

	"github.com/pkg/errors"
	"github.com/quic-go/quic-go"
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
	stream connection.ReadWriteStream
}

// New constructs a new publisher client.
func New() Client {
	return &client{}
}

func (c *client) Start(port int, connectionClosed chan struct{}) error {
	var err error
	c.stream, err = c.setupReadWriteStream(port)
	if err != nil {
		return errors.Wrap(err, "setup read write stream")
	}

	c.stream.SetConnClosedCallback(func() {
		connectionClosed <- struct{}{}
	})

	log.Tracef("Started listening to messages on port %d", port)
	return nil
}

func (c *client) setupReadWriteStream(port int) (connection.ReadWriteStream, error) {
	log.Trace("Setting up UDP connection...")
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: 0})
	if err != nil {
		return nil, errors.Wrap(err, "listen udp")
	}
	transport := &quic.Transport{Conn: udpConn}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	address := fmt.Sprintf("localhost:%d", port)
	log.Tracef("Dialing the server on address %q...", address)
	conn, err := transport.Dial(
		ctx, &net.UDPAddr{Port: port}, &tls.Config{InsecureSkipVerify: true}, &quic.Config{
			MaxIdleTimeout: DefaultIdleTimeout,
		})
	if err != nil {
		return nil, errors.Wrapf(err, "dial %q", address)
	}

	log.Trace("Accepting stream...")
	return connection.New(conn).
		AcceptReadWriteStream(ctx, c.handleMessage)
}

func (c *client) Publish(message string) error {
	if err := c.stream.SendMessage(entity.Message{
		Text: message,
	}); err != nil {
		return errors.Wrap(err, "send message")
	}

	log.Infof("Message %q successfully published", message)
	return nil
}

func (c *client) handleMessage(message entity.Message) {
	log.Infof("Received message: %q", message.Text)
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
