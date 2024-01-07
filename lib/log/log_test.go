package log

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFormatMessage(t *testing.T) {
	now = func() time.Time {
		return time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	message := "Hello, World!"
	got := formatMessage(message)
	require.Equal(t, "[2020-01-01 00:00:00] Hello, World!\n", got)
}
