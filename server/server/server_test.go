package server

import (
	"crypto/tls"
	"testing"

	"assignment/server/server/listener"
	listenermocks "assignment/server/server/listener/mocks"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestServer_Lifecycle(t *testing.T) {
	var (
		config = Config{
			SubscriberPort: 1111,
			PublisherPort:  2222,
			TLS:            &tls.Config{},
		}
		s            = New(config, zap.NewNop())
		ctrl         = gomock.NewController(t)
		listenerMock = listenermocks.NewMockListener(ctrl)
	)

	s.(*server).newListener = func(cb listener.NewConnectionCallback, _ *zap.Logger) listener.Listener {
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
				s            = New(config, zap.NewNop()).(*server)
			)

			tc.setup(listenerMock)
			s.newListener = func(cb listener.NewConnectionCallback, _ *zap.Logger) listener.Listener {
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
