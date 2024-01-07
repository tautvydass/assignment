package connection

import (
	"fmt"
	"sync"

	"assignment/lib/entity"

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
	CloseStream()
}

type readStream struct {
	sync.RWMutex
	messageReceiver MessageReceiver
	readBufferSize  int
	buffer          []byte

	stream quic.ReceiveStream
}

func NewReadStream(
	stream quic.ReceiveStream,
	messageReceiver MessageReceiver,
) ReadStream {
	rs := &readStream{
		messageReceiver: messageReceiver,
		readBufferSize:  DefaultReadBufferSize,
		stream:          stream,
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

func (s *readStream) CloseStream() {
	s.stream.CancelRead(0)
}

func (s *readStream) listen() {
	for {
		s.syncBuffer()
		size, err := s.stream.Read(s.buffer)
		if err != nil {
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
