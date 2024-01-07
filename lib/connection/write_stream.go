package connection

import (
	"time"

	"assignment/lib/entity"

	"github.com/pkg/errors"
	"github.com/quic-go/quic-go"
)

// WriteStream provides functionality for writing messages to a stream.
type WriteStream interface {
	// SendMessage sends a message to the stream.
	SendMessage(message entity.Message) error
	// SetSendMessageTimeout sets the timeout for sending a message.
	SetSendMessageTimeout(timeout time.Duration)
	// CloseStream closes the stream.
	CloseStream() error
}

type writeStream struct {
	// TODO: create a stripped interface alias for quic.SendStream and
	// use it instead of quic.SendStream.
	stream  quic.SendStream
	timeout time.Duration
}

// NewWriteStream constructs a new write stream.
func NewWriteStream(stream quic.SendStream) WriteStream {
	return &writeStream{
		stream: stream,
	}
}

func (s *writeStream) SendMessage(message entity.Message) error {
	if s.timeout != 0 {
		deadline := time.Now().Add(s.timeout)
		s.stream.SetWriteDeadline(deadline)
	}

	_, err := s.stream.Write(message.Bytes())
	if err != nil {
		return errors.Wrap(err, "write message")
	}
	return nil
}

func (s *writeStream) SetSendMessageTimeout(timeout time.Duration) {
	s.timeout = timeout
}

func (s *writeStream) CloseStream() error {
	return s.stream.Close()
}
