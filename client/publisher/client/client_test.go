package client

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

func TestClient(t *testing.T) {
	tlsConfig, err := certificate.LoadTLSConfig(
		"../../../testdata/test_server.crt", "../../../testdata/test_server.key")
	require.NoError(t, err)

	var serverSentMessage sync.WaitGroup
	var publisherSentMessage sync.WaitGroup
	serverSentMessage.Add(1)
	publisherSentMessage.Add(1)

	// Set up the server.
	listener, err := connection.StartListener(8085, tlsConfig)
	require.NoError(t, err)
	serverMessageCollector := testutil.NewMessageCollector()
	go func() {
		// Start the listener and wait until publisher client connects.
		serverConn, err := listener.Accept(context.Background())
		require.NoError(t, err)

		serverStream, err := connection.New(serverConn).OpenReadWriteStream(
			context.Background(), serverMessageCollector.Add)
		require.NoError(t, err)

		// TODO: do not use time.Sleep() in tests, find a better way
		time.Sleep(time.Millisecond * 500)

		// Send message to the publisher and wait until it sends a message back.
		require.NoError(t, serverStream.SendMessage(
			entity.Message{Text: "Hello from Server!"}))
		serverSentMessage.Done()
		publisherSentMessage.Wait()

		require.NoError(t, serverStream.CloseStream())
	}()

	// Set up the publisher client.
	client := New()
	connectionClosedCh := make(chan struct{})
	require.NoError(t, client.Start(8085, connectionClosedCh))
	publisherMessageCollector := testutil.NewMessageCollector()
	client.SetMessageReceiver(publisherMessageCollector.Add)

	// Wait until until server sends a message and send a message back.
	serverSentMessage.Wait()
	require.NoError(t, client.Publish("Hello from Publisher!"))

	// TODO: do not use time.Sleep() in tests, find a better way
	time.Sleep(time.Millisecond * 200)
	publisherSentMessage.Done()

	// Wait until server shuts down.
	<-connectionClosedCh
	require.NoError(t, client.Close())

	// Make sure that server and publisher exchanged messages.
	require.Equal(t, []entity.Message{
		{Text: "Hello from Server!"},
	}, publisherMessageCollector.Get())
	require.Equal(t, []entity.Message{
		{Text: "Hello from Publisher!"},
	}, serverMessageCollector.Get())
}
