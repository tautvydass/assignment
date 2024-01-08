package connection

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"assignment/lib/log"

	"github.com/pkg/errors"
	"github.com/quic-go/quic-go"
)

// TODO: add tests by creating interface type aliases
// for quic.Connection, quic.Stream, quic.ReceiveStream,
// and quic.SendStream, then generating mocks using those
// aliases. Currently, only the happy paths of this package
// are covered.

// DefaultIdleTimeout is the default idle timeout for the
// connection.
const DefaultIdleTimeout = time.Hour

// Connection is an interface for the connection.
type Connection interface {
	// OpenWriteStream opens a new unidirectional stream
	// for sending messages.
	OpenWriteStream(
		ctx context.Context,
	) (WriteStream, error)
	// AcceptReadStream accepts a new unidirectional stream
	//  for receiving messages.
	AcceptReadStream(
		ctx context.Context,
		messageReceiver MessageReceiver,
	) (ReadStream, error)
	// OpenReadWriteStream opens a new bidirectional stream
	// for sending and receiving messages.
	OpenReadWriteStream(
		ctx context.Context,
		messageReceiver MessageReceiver,
	) (ReadWriteStream, error)
	// AcceptReadWriteStream accepts a new bidirectional
	// stream for sending and receiving messages.
	AcceptReadWriteStream(
		ctx context.Context,
		messageReceiver MessageReceiver,
	) (ReadWriteStream, error)
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

	return NewWriteStream(c.conn, str), nil
}

func (c *connection) AcceptReadStream(
	ctx context.Context,
	messageReceiver MessageReceiver,
) (ReadStream, error) {
	str, err := c.conn.AcceptUniStream(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "accept unidirectional stream")
	}

	return NewReadStream(c.conn, str, messageReceiver), nil
}

func (c *connection) OpenReadWriteStream(
	ctx context.Context,
	messageReceiver MessageReceiver,
) (ReadWriteStream, error) {
	str, err := c.conn.OpenStreamSync(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "open stream")
	}

	return NewReadWriteStream(c.conn, str, messageReceiver), nil
}

func (c *connection) AcceptReadWriteStream(
	ctx context.Context,
	messageReceiver MessageReceiver,
) (ReadWriteStream, error) {
	str, err := c.conn.AcceptStream(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "accept stream")
	}

	return NewReadWriteStream(c.conn, str, messageReceiver), nil
}

// Connect dials the server on the given port and returns the connection.
func Connect(ctx context.Context, port int) (Connection, error) {
	// Set up UDP connection.
	log.Trace("Setting up UDP connection...")
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: 0})
	if err != nil {
		return nil, errors.Wrap(err, "listen udp")
	}

	// Set up QUIC transport.
	transport := &quic.Transport{Conn: udpConn}

	// Dial the server.
	address := fmt.Sprintf("localhost:%d", port)
	log.Tracef("Dialing the server on address %q...", address)
	conn, err := transport.Dial(
		ctx, &net.UDPAddr{Port: port}, &tls.Config{InsecureSkipVerify: true}, &quic.Config{
			MaxIdleTimeout: DefaultIdleTimeout,
		})
	if err != nil {
		return nil, errors.Wrapf(err, "dial %q", address)
	}

	return New(conn), nil
}
