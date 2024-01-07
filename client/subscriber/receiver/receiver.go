package receiver

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
	DefaultTimeout = time.Hour
	// DefaultIdleTimeout is the default idle timeout for the
	// connection.
	DefaultIdleTimeout = time.Hour
)

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
	r.readStream, err = r.setupReadStream(port, connectionClosed)
	if err != nil {
		return errors.Wrap(err, "setup read stream")
	}

	log.Infof("Started listening to messages on port %d", port)
	return nil
}

func (r *receiver) setupReadStream(
	port int, connectionClosed chan struct{},
) (connection.ReadStream, error) {
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

	log.Trace("Accepting read stream and waiting for messages...")
	return connection.New(conn).AcceptReadStream(ctx, r.handleMessage, connectionClosed)
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
