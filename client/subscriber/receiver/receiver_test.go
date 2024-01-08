package receiver

import (
	"context"
	"sync"
	"testing"
	"time"

	"assignment/lib/certificate"
	"assignment/lib/connection"
	"assignment/lib/entity"
	"assignment/lib/testutil"

	"github.com/stretchr/testify/require"
)

func TestReceiver(t *testing.T) {
	tlsConfig, err := certificate.LoadTLSConfig(
		"../../../testdata/test_server.crt", "../../../testdata/test_server.key")
	require.NoError(t, err)

	var serverMessageSent sync.WaitGroup
	var subscriberReceivedMessage sync.WaitGroup
	serverMessageSent.Add(1)
	subscriberReceivedMessage.Add(1)

	// Set up the server.
	listener, err := connection.StartListener(8086, tlsConfig)
	require.NoError(t, err)

	go func() {
		// Start the listener and wait until subscriber client connects.
		serverConn, err := listener.Accept(context.Background())
		require.NoError(t, err)

		serverStream, err := connection.New(serverConn).OpenWriteStream(
			context.Background())
		require.NoError(t, err)

		// TODO: do not use time.Sleep() in tests, find a better way
		time.Sleep(time.Millisecond * 500)

		// Send message to the subscriver and wait until it receives it.
		require.NoError(t, serverStream.SendMessage(
			entity.Message{Text: "Hello from Server!"}))
		serverMessageSent.Done()
		subscriberReceivedMessage.Wait()

		require.NoError(t, serverStream.CloseStream())
	}()

	receiver := New()
	connectionClosedCh := make(chan struct{})
	require.NoError(t, receiver.Start(8086, connectionClosedCh))
	subscriberMessageCollector := testutil.NewMessageCollector()
	receiver.SetMessageReceiver(func(message entity.Message) {
		subscriberMessageCollector.Add(message)
		subscriberReceivedMessage.Done()
	})

	// Wait until server sends the message and shuts down.
	<-connectionClosedCh
	require.NoError(t, receiver.Close())

	// Make sure the subscriber received the message.
	require.Equal(t, []entity.Message{
		{Text: "Hello from Server!"},
	}, subscriberMessageCollector.Get())
}
