package client

import (
	"context"
	"time"

	"assignment/lib/connection"
	"assignment/lib/entity"
	"assignment/lib/log"

	"github.com/pkg/errors"
)

// DefaultTimeout is the default timeout for establishing
// a connection to the servec.
const DefaultTimeout = time.Hour

// Client represents subscriber client that receives messages
// from the servec. The receiver will log all received messages.
type Client interface {
	// Start establishes a connection with the server and
	// begins listening to messages. Given channel is closed
	// when connection is closed by the servec.
	Start(port int, connectionClosed chan struct{}) error
	// SetMessageReceiver sets the message receiver callback.
	SetMessageReceiver(receiver connection.MessageReceiver)
	// Close closes the connection with the servec.
	Close() error
}

type client struct {
	readStream connection.ReadStream
}

// New constructs a new subscriber client.
func New() Client {
	return &client{}
}

func (c *client) Start(port int, connectionClosed chan struct{}) error {
	var err error
	c.readStream, err = c.setupReadStream(port)
	if err != nil {
		return errors.Wrap(err, "setup read stream")
	}

	c.readStream.SetConnClosedCallback(func() {
		connectionClosed <- struct{}{}
	})

	log.Infof("Started listening to messages on port %d", port)
	return nil
}

func (c *client) SetMessageReceiver(receiver connection.MessageReceiver) {
	c.readStream.SetMessageReceiver(receiver)
}

func (c *client) setupReadStream(port int) (connection.ReadStream, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	conn, err := connection.Connect(ctx, port)
	if err != nil {
		return nil, errors.Wrap(err, "connect")
	}

	log.Trace("Accepting read stream and waiting for messages...")
	return conn.AcceptReadStream(ctx, c.handleMessage)
}

func (c *client) handleMessage(message entity.Message) {
	log.Infof("Received message: %q", message.Text)
}

func (c *client) Close() error {
	if c.readStream == nil {
		return nil
	}
	return errors.Wrap(
		c.readStream.CloseStream(),
		"close read stream",
	)
}
