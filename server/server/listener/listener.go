package listener

import (
	"context"
	"crypto/tls"

	"assignment/lib/connection"
	"assignment/lib/log"

	"github.com/pkg/errors"
)

// ErrAlreadyStarted is returned when attempting to start a
// listener that's already started.
var ErrAlreadyStarted = errors.New("listener already started")

// NewConnectionCallback is the type alias for new connection
// callback function.
type NewConnectionCallback func(conn connection.Connection)

// Listener is an interface for the connection listener.
type Listener interface {
	// Start starts the listener on the given port on a separate goroutine.
	Start(port int, tlsConfig *tls.Config) error
	// Shutdown shuts down the listener.
	Shutdown() error
}

type listener struct {
	callback     NewConnectionCallback
	close        chan struct{}
	listener     connection.QUICListener
	started      bool
	connCancelFn context.CancelFunc

	// used for mocks in tests
	startListenerFn func(port int, tlsConfig *tls.Config) (connection.QUICListener, error)
}

// New creates a new connection listener. Provided callback function
// must be goroutine safe.
func New(cb NewConnectionCallback) Listener {
	return &listener{
		callback:        cb,
		startListenerFn: connection.StartListener,
	}
}

func (l *listener) Start(port int, tlsConfig *tls.Config) error {
	if l.started {
		return ErrAlreadyStarted
	}

	listener, err := l.startListenerFn(port, tlsConfig)
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
	ctx, cancel := context.WithCancel(context.Background())
	l.connCancelFn = cancel

	for {
		select {
		case <-l.close:
			l.started = false
			return
		default:
			conn, err := l.listener.Accept(ctx)
			if err != nil {
				// Ignore context canceled error as it's expected only when
				// the server is shutting down. Log the error otherwise.
				if !errors.Is(err, context.Canceled) {
					log.Errorf("Error accepting connection: %v", err)
				}
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

	if l.connCancelFn != nil {
		l.connCancelFn()
	}
	l.close <- struct{}{}

	return l.listener.Close()
}
