package entity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessageConversion(t *testing.T) {
	message := Message{
		Text: "Hello, World!",
	}
	require.Equal(t, []byte("Hello, World!"), message.Bytes())

	bytes := []byte("Hello, World!")
	require.Equal(t, message, MessageFromBytes(bytes))
}
