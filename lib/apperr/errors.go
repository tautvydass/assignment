package apperr

import (
	"strings"

	"github.com/quic-go/quic-go"
)

// canceledStreamMessage is the partial error message for a canceled stream.
const canceledStreamMessage = "close called for canceled stream"

// IsConnectionClosedByPeerErr returns true if the given error is caused
// by the other peer terminating the connection.
func IsConnectionClosedByPeerErr(err error) bool {
	if streamErr, ok := err.(*quic.StreamError); ok {
		return streamErr.ErrorCode == ErrCodeClosedByClient
	}
	// TODO: find a better way to check for this error.
	if strings.Contains(err.Error(), canceledStreamMessage) {
		return true
	}
	if appErr, ok := err.(*quic.ApplicationError); ok {
		return appErr.ErrorCode == ErrCodeClosedByClient
	}
	return false
}
