package connection

import (
	"time"

	"assignment/lib/apperr"
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
	conn    quic.Connection
	stream  quic.SendStream
	timeout time.Duration
}

// NewWriteStream constructs a new write stream.
// TODO: consider implementing a ping mechanism to check if the
// connection is still alive.
func NewWriteStream(conn quic.Connection, stream quic.SendStream) WriteStream {
	return &writeStream{
		conn:   conn,
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
	s.stream.CancelWrite(apperr.ErrCodeClosedByClient)
	if err := s.stream.Close(); err != nil {
		if !apperr.IsConnectionClosedByPeerErr(err) {
			return errors.Wrap(err, "close write stream")
		}
	}

	return errors.Wrap(
		s.conn.CloseWithError(apperr.ErrCodeClosedByClient, ""),
		"close connection")
}
