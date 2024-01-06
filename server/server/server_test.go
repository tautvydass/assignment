package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServer_Lifecycle(t *testing.T) {
	var (
		config = Config{
			SubscriberPort: 1111,
			PublisherPort:  2222,
		}
		s = New(config)
	)

	// make sure config is set
	require.Equal(t, config, s.(*server).config)

	// make sure server has started
	require.NoError(t, s.Start())
	require.True(t, s.(*server).started)

	// make sure server can't be started again
	require.EqualError(t, s.Start(), ErrAlreadyStarted.Error())

	// should not return any error
	require.NoError(t, s.Shutdown())
	require.False(t, s.(*server).started)
}
