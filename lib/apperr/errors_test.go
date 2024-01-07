package apperr

import (
	"testing"

	"github.com/quic-go/quic-go"
	"github.com/stretchr/testify/assert"
)

func TestIsConnectionClosedByPeerErr(t *testing.T) {
	tests := map[string]struct {
		err  error
		want bool
	}{
		"stream_closed_error": {
			err: &quic.StreamError{
				ErrorCode: ErrCodeClosedByClient,
			},
			want: true,
		},
		"other_stream_error": {
			err: &quic.StreamError{
				ErrorCode: 999999,
			},
			want: false,
		},
		"stream_cancelled_error": {
			err:  errCancelledStream,
			want: true,
		},
		"closed_by_application_error": {
			err: &quic.ApplicationError{
				ErrorCode: ErrCodeClosedByClient,
			},
			want: true,
		},
		"other_application_error": {
			err: &quic.ApplicationError{
				ErrorCode: 999999,
			},
			want: false,
		},
		"other_error": {
			err:  assert.AnError,
			want: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := IsConnectionClosedByPeerErr(tc.err)
			assert.Equal(t, tc.want, got)
		})
	}
}
