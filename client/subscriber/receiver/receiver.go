package receiver

import (
	"context"
	"time"

	"assignment/lib/connection"
	"assignment/lib/entity"
	"assignment/lib/log"

	"github.com/pkg/errors"
)

// DefaultTimeout is the default timeout for establishing
// a connection to the server.
const DefaultTimeout = time.Hour

// Receiver is an interface for receiving messages from
// the server. The receiver will log all received messages.
type Receiver interface {
	// Start establishes a connection with the server and
	// begins listening to messages. Given channel is closed
	// when connection is closed by the server.
	Start(port int, connectionClosed chan struct{}) error
	// Close closes the connection with the server.
	Close() error
}

type receiver struct {
	readStream connection.ReadStream
}

// New constructs a new receiver.
func New() Receiver {
	return &receiver{}
}

func (r *receiver) Start(port int, connectionClosed chan struct{}) error {
	var err error
	r.readStream, err = r.setupReadStream(port)
	if err != nil {
		return errors.Wrap(err, "setup read stream")
	}

	r.readStream.SetConnClosedCallback(func() {
		connectionClosed <- struct{}{}
	})

	log.Infof("Started listening to messages on port %d", port)
	return nil
}

func (r *receiver) setupReadStream(port int) (connection.ReadStream, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	conn, err := connection.Connect(ctx, port)
	if err != nil {
		return nil, errors.Wrap(err, "connect")
	}

	log.Trace("Accepting read stream and waiting for messages...")
	return conn.AcceptReadStream(ctx, r.handleMessage)
}

func (r *receiver) handleMessage(message entity.Message) {
	log.Infof("Received message: %q", message.Text)
}

func (r *receiver) Close() error {
	if r.readStream == nil {
		return nil
	}
	return errors.Wrap(
		r.readStream.CloseStream(),
		"close read stream",
	)
}
