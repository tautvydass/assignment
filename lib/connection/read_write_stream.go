package connection

import (
	"time"

	"assignment/lib/apperr"
	"assignment/lib/entity"

	"github.com/pkg/errors"
	"github.com/quic-go/quic-go"
)

// ReadWriteStream provides functionality for reading messages
// from the server and sending messages to the server.
type ReadWriteStream interface {
	ReadStream
	WriteStream
}

type readWriteStream struct {
	conn   quic.Connection
	stream quic.Stream

	readStream  ReadStream
	writeStream WriteStream
}

// NewReadWriteStream constructs a new read write stream.
func NewReadWriteStream(
	conn quic.Connection,
	stream quic.Stream,
	messageReceiver MessageReceiver,
) ReadWriteStream {
	return &readWriteStream{
		conn:   conn,
		stream: stream,

		readStream:  NewReadStream(conn, stream, messageReceiver),
		writeStream: NewWriteStream(conn, stream),
	}
}

func (s *readWriteStream) SetMessageReceiver(messageReceiver MessageReceiver) {
	s.readStream.SetMessageReceiver(messageReceiver)
}

func (s *readWriteStream) SetConnClosedCallback(connClosedCallback ConnClosedCallback) {
	s.readStream.SetConnClosedCallback(connClosedCallback)
}

func (s *readWriteStream) SetReadBufferSize(size int) {
	s.readStream.SetReadBufferSize(size)
}

func (s *readWriteStream) SendMessage(message entity.Message) error {
	return s.writeStream.SendMessage(message)
}

func (s *readWriteStream) SetSendMessageTimeout(timeout time.Duration) {
	s.writeStream.SetSendMessageTimeout(timeout)
}

func (s *readWriteStream) CloseStream() error {
	s.stream.CancelRead(apperr.ErrCodeClosedByClient)
	s.stream.CancelWrite(apperr.ErrCodeClosedByClient)

	if err := s.stream.Close(); err != nil {
		if !apperr.IsConnectionClosedByPeerErr(err) {
			return errors.Wrap(err, "close stream")
		}
	}

	return errors.Wrap(
		s.conn.CloseWithError(apperr.ErrCodeClosedByClient, ""),
		"close connection")
}
