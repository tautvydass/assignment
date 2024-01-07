package connection

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWriteStream_SetSendMessageTimeout(t *testing.T) {
	str := &writeStream{}
	require.Empty(t, str.timeout)

	timeout := time.Minute
	str.SetSendMessageTimeout(timeout)
	require.Equal(t, timeout, str.timeout)
}
