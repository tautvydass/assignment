package server

import (
	"context"
	"crypto/tls"
	"sync"
	"testing"
	"time"

	"assignment/lib/certificate"
	"assignment/lib/connection"
	connectionmocks "assignment/lib/connection/mocks"
	"assignment/lib/entity"
	"assignment/server/server/controller"
	controllermocks "assignment/server/server/controller/mocks"
	"assignment/server/server/listener"
	listenermocks "assignment/server/server/listener/mocks"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	tlsConfig, err := certificate.LoadTLSConfig(
		"testdata/test_server.crt", "testdata/test_server.key")
	require.NoError(t, err)

	config := Config{
		SubscriberPort:     8083,
		PublisherPort:      8084,
		TLS:                tlsConfig,
		OpenStreamTimeout:  time.Second,
		SendMessageTimeout: time.Second,
	}
	server := New(config)

	require.NoError(t, server.Start())

	// Connect to the server as a publisher.
	publisherConn, err := connection.Connect(
		context.Background(), config.PublisherPort)
	require.NoError(t, err)
	publisherMessageCollector := newMessageCollector()
	publisherStream, err := publisherConn.AcceptReadWriteStream(
		context.Background(), publisherMessageCollector.add)
	require.NoError(t, err)

	// TODO: do not use time.Sleep() in tests, find a better way
	time.Sleep(time.Millisecond * 100)

	// Connect to the server as a subscriber.
	subscriberConn, err := connection.Connect(
		context.Background(), config.SubscriberPort)
	require.NoError(t, err)
	subscriberMessageCollector := newMessageCollector()
	subscriberStream, err := subscriberConn.AcceptReadStream(
		context.Background(), subscriberMessageCollector.add)
	require.NoError(t, err)

	// TODO: do not use time.Sleep() in tests, find a better way
	time.Sleep(time.Millisecond * 100)

	// Send a message from the publisher to the subscriber.
	publisherMessage := entity.Message{Text: "New message from publisher"}
	require.NoError(t, publisherStream.SendMessage(publisherMessage))

	// TODO: do not use time.Sleep() in tests, find a better way
	time.Sleep(time.Millisecond * 100)

	require.NoError(t, publisherStream.CloseStream())
	require.NoError(t, subscriberStream.CloseStream())
	require.NoError(t, server.Shutdown())

	// TODO: do not use time.Sleep() in tests, find a better way
	time.Sleep(time.Millisecond * 100)

	// Make sure that the subscriber has received the messages.
	require.Equal(t, []entity.Message{
		{Text: controller.MessageHelloSubscriber},
		publisherMessage,
	}, subscriberMessageCollector.get())

	// Make sure that the publisher has received the messages.
	require.Equal(t, []entity.Message{
		{Text: controller.MessageNoSubscribers},
		{Text: controller.MessageNewSubscriber},
	}, publisherMessageCollector.get())
}

func TestServer_Lifecycle(t *testing.T) {
	var (
		config = Config{
			SubscriberPort: 1111,
			PublisherPort:  2222,
			TLS:            &tls.Config{},
		}
		s            = New(config)
		ctrl         = gomock.NewController(t)
		listenerMock = listenermocks.NewMockListener(ctrl)
	)

	s.(*server).newListener = func(cb listener.NewConnectionCallback) listener.Listener {
		return listenerMock
	}

	listenerMock.EXPECT().Start(config.PublisherPort, gomock.Any()).Return(nil).Times(1)
	listenerMock.EXPECT().Start(config.SubscriberPort, gomock.Any()).Return(nil).Times(1)
	listenerMock.EXPECT().Shutdown().Return(nil).Times(2)

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

func TestServer_Start(t *testing.T) {
	var (
		config = Config{
			PublisherPort:  1111,
			SubscriberPort: 2222,
			TLS:            &tls.Config{},
		}
		tests = map[string]struct {
			setup   func(lm *listenermocks.MockListener)
			wantErr error
		}{
			"error_starting_publisher_listener": {
				setup: func(lm *listenermocks.MockListener) {
					lm.EXPECT().Start(config.PublisherPort, gomock.Any()).
						Return(assert.AnError).Times(1)
				},
				wantErr: errors.Wrap(assert.AnError, "start publisher listener"),
			},
			"error_starting_subscriber_listener": {
				setup: func(lm *listenermocks.MockListener) {
					lm.EXPECT().Start(config.PublisherPort, gomock.Any()).
						Return(nil).Times(1)
					lm.EXPECT().Start(config.SubscriberPort, gomock.Any()).
						Return(assert.AnError).Times(1)
				},
				wantErr: errors.Wrap(assert.AnError, "start subscriber listener"),
			},
			"happy_path": {
				setup: func(lm *listenermocks.MockListener) {
					lm.EXPECT().Start(config.PublisherPort, gomock.Any()).
						Return(nil).Times(1)
					lm.EXPECT().Start(config.SubscriberPort, gomock.Any()).
						Return(nil).Times(1)
				},
				wantErr: nil,
			},
		}
	)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				ctrl         = gomock.NewController(t)
				listenerMock = listenermocks.NewMockListener(ctrl)
				s            = New(config).(*server)
			)

			tc.setup(listenerMock)
			s.newListener = func(cb listener.NewConnectionCallback) listener.Listener {
				return listenerMock
			}

			err := s.Start()
			if tc.wantErr != nil {
				require.EqualError(t, err, tc.wantErr.Error())
				assert.False(t, s.started)
				return
			}

			require.NoError(t, err)
			assert.True(t, s.started)
		})
	}
}

func TestServer_addPublisher(t *testing.T) {
	type mocks struct {
		conn       *connectionmocks.MockConnection
		stream     *connectionmocks.MockReadWriteStream
		controller *controllermocks.MockCommsController
	}

	var (
		config = Config{
			SendMessageTimeout: time.Second,
			OpenStreamTimeout:  time.Minute,
		}
		messageReceiver = func(_ entity.Message) {}
		tests           = map[string]struct {
			setup func(m mocks)
		}{
			"error_opening_stream": {
				setup: func(m mocks) {
					m.controller.EXPECT().MessageReceiver().
						Return(messageReceiver).Times(1)
					m.conn.EXPECT().OpenReadWriteStream(gomock.Any(), gomock.Any()).
						Return(nil, assert.AnError).Times(1)
				},
			},
			"happy_path": {
				setup: func(m mocks) {
					m.controller.EXPECT().MessageReceiver().
						Return(messageReceiver).Times(1)
					m.conn.EXPECT().OpenReadWriteStream(gomock.Any(), gomock.Any()).
						DoAndReturn(
							func(_ context.Context, mr connection.MessageReceiver) (connection.ReadWriteStream, error) {
								require.NotNil(t, mr)
								return m.stream, nil
							}).Times(1)
					m.stream.EXPECT().SetSendMessageTimeout(config.SendMessageTimeout).Times(1)
					m.controller.EXPECT().AddPublisher(m.stream).Times(1)
				},
			},
		}
	)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				ctrl           = gomock.NewController(t)
				connectionMock = connectionmocks.NewMockConnection(ctrl)
				streamMock     = connectionmocks.NewMockReadWriteStream(ctrl)
				controllerMock = controllermocks.NewMockCommsController(ctrl)
			)

			tc.setup(mocks{
				conn:       connectionMock,
				stream:     streamMock,
				controller: controllerMock,
			})
			s := &server{
				config:          config,
				commsController: controllerMock,
			}

			s.addPublisher(connectionMock)
		})
	}
}

func TestServer_addSubscriber(t *testing.T) {
	type mocks struct {
		conn       *connectionmocks.MockConnection
		stream     *connectionmocks.MockReadWriteStream
		controller *controllermocks.MockCommsController
	}

	var (
		config = Config{
			SendMessageTimeout: time.Second,
			OpenStreamTimeout:  time.Minute,
		}
		tests = map[string]struct {
			setup func(m mocks)
		}{
			"error_opening_stream": {
				setup: func(m mocks) {
					m.conn.EXPECT().OpenWriteStream(gomock.Any()).
						Return(nil, assert.AnError).Times(1)
				},
			},
			"happy_path": {
				setup: func(m mocks) {
					m.conn.EXPECT().OpenWriteStream(gomock.Any()).
						Return(m.stream, nil).Times(1)
					m.stream.EXPECT().SetSendMessageTimeout(config.SendMessageTimeout).Times(1)
					m.controller.EXPECT().AddSubscriber(m.stream).Times(1)
				},
			},
		}
	)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				ctrl           = gomock.NewController(t)
				connectionMock = connectionmocks.NewMockConnection(ctrl)
				streamMock     = connectionmocks.NewMockReadWriteStream(ctrl)
				controllerMock = controllermocks.NewMockCommsController(ctrl)
			)

			tc.setup(mocks{
				conn:       connectionMock,
				stream:     streamMock,
				controller: controllerMock,
			})
			s := &server{
				config:          config,
				commsController: controllerMock,
			}

			s.addSubscriber(connectionMock)
		})
	}
}

type messageCollector struct {
	sync.RWMutex

	messages []entity.Message
}

func (m *messageCollector) add(msg entity.Message) {
	m.Lock()
	defer m.Unlock()

	m.messages = append(m.messages, msg)
}

func (m *messageCollector) get() []entity.Message {
	m.RLock()
	defer m.RUnlock()

	return m.messages
}

func newMessageCollector() *messageCollector {
	return &messageCollector{
		messages: make([]entity.Message, 0),
	}
}
