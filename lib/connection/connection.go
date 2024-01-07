package connection

import (
	"context"

	"github.com/pkg/errors"
	"github.com/quic-go/quic-go"
)

// Connection is an interface for the connection.
type Connection interface {
	// OpenWriteStream opens a new unidirectional stream for sending messages.
	OpenWriteStream(
		ctx context.Context,
	) (WriteStream, error)
	// AcceptReadStream accepts a new unidirectional stream for receiving messages.
	AcceptReadStream(
		ctx context.Context,
		messageReceiver MessageReceiver,
		connectionClosed chan struct{},
	) (ReadStream, error)
}

type connection struct {
	conn quic.Connection
}

// New constructs a new connection.
func New(conn quic.Connection) Connection {
	return &connection{
		conn: conn,
	}
}

func (c *connection) OpenWriteStream(ctx context.Context) (WriteStream, error) {
	str, err := c.conn.OpenUniStreamSync(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "open unidirectional stream")
	}

	return NewWriteStream(str), nil
}

func (c *connection) AcceptReadStream(
	ctx context.Context,
	messageReceiver MessageReceiver,
	connectionClosed chan struct{},
) (ReadStream, error) {
	str, err := c.conn.AcceptUniStream(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "accept unidirectional stream")
	}

	return NewReadStream(str, messageReceiver, connectionClosed), nil
}
