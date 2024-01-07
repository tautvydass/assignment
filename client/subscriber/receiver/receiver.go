package receiver

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

// DefaultTimeout is the default timeout for establishing
// a connection to the server
const DefaultTimeout = time.Second * 30

// Receiver is an interface for receiving messages from
// the server.
type Receiver interface {
	// Start establishes a connection with the server and
	// begins listening to messages. Given channel is closed
	// when connection is closed by the server.
	Start(port int, connectionClosed chan struct{}) error
	// Close closes the connection with the server.
	Close() error
}

type receiver struct {
	logger     *zap.Logger
	readStream connection.ReadStream
}

// New constructs a new receiver.
func New(logger *zap.Logger) Receiver {
	return &receiver{
		logger: logger,
	}
}

func (r *receiver) Start(port int, connectionClosed chan struct{}) error {
	var err error
	r.readStream, err = r.setupReadStream(port, connectionClosed)
	if err != nil {
		return errors.Wrap(err, "setup read stream")
	}

	r.logger.Info("Started listening to messages", zap.Int("port", port))
	return nil
}

func (r *receiver) setupReadStream(
	port int, connectionClosed chan struct{},
) (connection.ReadStream, error) {
	r.logger.Info("Setting up UDP connection")
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: 0})
	if err != nil {
		return nil, errors.Wrap(err, "listen udp")
	}
	transport := &quic.Transport{Conn: udpConn}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	address := fmt.Sprintf("localhost:%d", port)
	r.logger.Info("Dialing the server", zap.String("address", address))
	conn, err := transport.Dial(
		ctx, &net.UDPAddr{Port: port}, &tls.Config{InsecureSkipVerify: true}, &quic.Config{})
	if err != nil {
		return nil, errors.Wrapf(err, "dial %q", address)
	}

	r.logger.Info("Accepting uni directional read stream")
	stream, err := conn.AcceptUniStream(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "accept unidirectional stream")
	}

	return connection.NewReadStream(conn, stream, r.handleMessage, connectionClosed), nil
}

func (r *receiver) handleMessage(message entity.Message) {
	r.logger.Info("Received message", zap.String("message", message.Text))
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
