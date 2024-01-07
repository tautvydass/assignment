package connection

import (
	"fmt"
	"io"
	"sync"

	"assignment/lib/apperr"
	"assignment/lib/entity"

	"github.com/pkg/errors"
	"github.com/quic-go/quic-go"
)

// DefaultReadBufferSize is the default size of the read buffer.
const DefaultReadBufferSize = 1024 * 1024

// MessageReceiver is the function callback for receiving messages.
type MessageReceiver func(message entity.Message)

// ReadStream provides functionality for reading messages
// from a stream.
type ReadStream interface {
	// SetMessageReceiver sets the message receiver.
	SetMessageReceiver(messageReceiver MessageReceiver)
	// SetReadBufferSize sets the read buffer size.
	SetReadBufferSize(size int)
	// CloseStream closes the stream.
	CloseStream() error
}

type readStream struct {
	sync.RWMutex
	messageReceiver MessageReceiver
	readBufferSize  int
	buffer          []byte

	stream           quic.ReceiveStream
	conn             quic.Connection
	connectionClosed chan struct{}
}

// NewReadStream constructs a new read stream.
func NewReadStream(
	conn quic.Connection,
	stream quic.ReceiveStream,
	messageReceiver MessageReceiver,
	connectionClosed chan struct{},
) ReadStream {
	rs := &readStream{
		messageReceiver:  messageReceiver,
		readBufferSize:   DefaultReadBufferSize,
		stream:           stream,
		conn:             conn,
		connectionClosed: connectionClosed,
	}

	go rs.listen()
	return rs
}

func (s *readStream) SetMessageReceiver(messageReceiver MessageReceiver) {
	s.Lock()
	defer s.Unlock()
	s.messageReceiver = messageReceiver
}

func (s *readStream) SetReadBufferSize(size int) {
	s.Lock()
	defer s.Unlock()
	s.readBufferSize = size
}

func (s *readStream) CloseStream() error {
	s.stream.CancelRead(0)
	return errors.Wrap(
		s.conn.CloseWithError(apperr.ErrCodeClosedByClient, ""),
		"close connection with error",
	)

}

func (s *readStream) listen() {
	for {
		s.syncBuffer()
		size, err := s.stream.Read(s.buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				// Connection closed by the server.
				close(s.connectionClosed)
				return
			}

			fmt.Printf("read stream error: %v\n", err)
			return
		}

		s.handleMessage(size)
	}
}

func (s *readStream) syncBuffer() {
	s.Lock()
	defer s.Unlock()
	if len(s.buffer) != s.readBufferSize {
		s.buffer = make([]byte, s.readBufferSize)
	}
}

func (s *readStream) handleMessage(messageSize int) {
	s.RLock()
	defer s.RUnlock()

	message := entity.MessageFromBytes(s.buffer[:messageSize])
	go s.messageReceiver(message)
}
