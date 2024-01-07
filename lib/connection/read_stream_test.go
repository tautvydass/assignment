package connection

import (
	"testing"

	"assignment/lib/entity"

	"github.com/stretchr/testify/require"
)

func TestReadStream_SetMessageReceiver(t *testing.T) {
	str := &readStream{}
	require.Nil(t, str.messageReceiver)

	called := 0
	messageReceiver := func(message entity.Message) {
		called++
	}
	str.SetMessageReceiver(messageReceiver)

	messageReceiver(entity.Message{})
	str.messageReceiver(entity.Message{})

	require.Equal(t, 2, called)
}

func TestReadStream_SetReadBufferSize(t *testing.T) {
	str := &readStream{}
	size := 1024
	str.SetReadBufferSize(size)
	require.Equal(t, size, str.readBufferSize)
}

func TestReadStream_syncBuffer(t *testing.T) {
	str := &readStream{
		readBufferSize: 128,
		buffer:         make([]byte, 0),
	}
	str.syncBuffer()
	require.Equal(t, 128, len(str.buffer))
}

func TestReadStream_handleMessage(t *testing.T) {
	str := &readStream{
		readBufferSize: 13,
		buffer:         []byte("Hello, World!"),
	}

	str.messageReceiver = func(message entity.Message) {
		require.Equal(t, "Hello, World!", message.Text)
	}
	str.handleMessage(str.readBufferSize)

	str.messageReceiver = func(message entity.Message) {
		require.Equal(t, "Hello", message.Text)
	}
	str.handleMessage(5)
}
