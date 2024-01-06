package listener

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"assignment/lib/connection"

	"github.com/pkg/errors"
	"github.com/quic-go/quic-go"
)

// ErrAlreadyStarted is returned when attempting to start a
// listener that's already started.
var ErrAlreadyStarted = errors.New("listener already started")

// NewConnectionCallback is the type alias for new connection
// callback function.
type NewConnectionCallback func(conn connection.Connection)

// Listener is an interface for the connection listener.
type Listener interface {
	Start(port int) error
	Shutdown() error
}

// QUICListener is an interface for the QUIC listener.
type QUICListener interface {
	Accept(context.Context) (quic.Connection, error)
	Close() error
}

type listener struct {
	callback NewConnectionCallback
	close    chan struct{}
	listener QUICListener
	started  bool

	// used for mocks in tests
	startListenerFn func(port int) (QUICListener, error)
}

// New creates a new connection listener. Provided callback function
// must be goroutine safe.
func New(cb NewConnectionCallback) Listener {
	return &listener{
		callback:        cb,
		startListenerFn: startListener,
	}
}

func (l *listener) Start(port int) error {
	if l.started {
		return ErrAlreadyStarted
	}

	listener, err := l.startListenerFn(port)
	if err != nil {
		return errors.Wrap(err, "start listener")
	}

	l.listener = listener
	l.close = make(chan struct{})
	l.started = true

	go l.run()
	return nil
}

func (l *listener) run() {
	for {
		select {
		case <-l.close:
			l.started = false
			return
		default:
			conn, err := l.listener.Accept(context.Background())
			if err != nil {
				fmt.Printf("error accepting connection: %v\n", err)
				continue
			}

			go l.callback(connection.New(conn))
		}
	}
}

func (l *listener) Shutdown() error {
	if !l.started {
		return nil
	}

	l.close <- struct{}{}
	return l.listener.Close()
}

func startListener(port int) (QUICListener, error) {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: port})
	if err != nil {
		return nil, errors.Wrap(err, "set up udp listener")
	}

	transport := &quic.Transport{
		Conn: conn,
	}

	// TODO: set up TLS certificate
	listener, err := transport.Listen(&tls.Config{}, &quic.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "set up quic listener")
	}

	return listener, nil
}
