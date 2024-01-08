package connection

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/pkg/errors"
	"github.com/quic-go/quic-go"
)

// QUICListener is an interface for the QUIC listener.
type QUICListener interface {
	// Accept accepts a new connection.
	Accept(context.Context) (quic.Connection, error)
	// Close closes the listener.
	Close() error
}

// StartListener starts a new QUIC listener on the given port.
func StartListener(port int, tlsConfig *tls.Config) (QUICListener, error) {
	// Set up UDP connection.
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: port})
	if err != nil {
		return nil, errors.Wrap(err, "set up udp listener")
	}

	// Set up QUIC transport.
	transport := &quic.Transport{Conn: udpConn}

	// Start the listener.
	listener, err := transport.Listen(tlsConfig, &quic.Config{
		MaxIdleTimeout: DefaultIdleTimeout,
	})
	if err != nil {
		return nil, errors.Wrap(err, "set up quic listener")
	}
	return listener, nil
}
