package connection

import (
	"testing"
	"time"

	"assignment/lib/entity"

	"github.com/stretchr/testify/require"
)

func TestReadWriteStream_SetMessageReceiver(t *testing.T) {
	var (
		message         = entity.Message{Text: "message"}
		messageReceiver = func(m entity.Message) {
			require.Equal(t, message, m)
		}
		readStream      = &readStream{}
		readWriteStream = &readWriteStream{
			readStream: readStream,
		}
	)

	readWriteStream.SetMessageReceiver(messageReceiver)
	readStream.messageReceiver(message)
}

func TestReadWriteStream_SetReadBufferSize(t *testing.T) {
	var (
		size            = 123
		readStream      = &readStream{}
		readWriteStream = &readWriteStream{
			readStream: readStream,
		}
	)

	readWriteStream.SetReadBufferSize(size)
	require.Equal(t, size, readStream.readBufferSize)
}

func TestReadWriteStream_SetSendMessageTimeout(t *testing.T) {
	var (
		timeout         = time.Minute
		writeStream     = &writeStream{}
		readWriteStream = &readWriteStream{
			writeStream: writeStream,
		}
	)

	readWriteStream.SetSendMessageTimeout(timeout)
	require.Equal(t, timeout, writeStream.timeout)
}
