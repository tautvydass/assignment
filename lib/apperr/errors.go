package apperr

import (
	"errors"

	"github.com/quic-go/quic-go"
)

// errCancelledStream is an untyped error returned when the stream is canceled.
var errCancelledStream = errors.New("close called for canceled stream 3")

// IsConnectionClosedByPeerErr returns true if the given error is caused
// by the other peer terminating the connection.
func IsConnectionClosedByPeerErr(err error) bool {
	if streamErr, ok := err.(*quic.StreamError); ok {
		return streamErr.ErrorCode == ErrCodeClosedByClient
	}
	// TODO: find a better way to check for this error.
	if err.Error() == errCancelledStream.Error() {
		return true
	}
	if appErr, ok := err.(*quic.ApplicationError); ok {
		return appErr.ErrorCode == ErrCodeClosedByClient
	}
	return false
}
