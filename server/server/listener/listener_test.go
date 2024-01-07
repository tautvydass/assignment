package listener

import (
	"crypto/tls"
	"sync"
	"testing"

	"assignment/lib/connection"
	"assignment/server/server/listener/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListener(t *testing.T) {
	var (
		ctrl            = gomock.NewController(t)
		listenerMock    = mocks.NewMockQUICListener(ctrl)
		startListenerFn = func(int, *tls.Config) (QUICListener, error) {
			return listenerMock, nil
		}
	)

	gomock.InOrder(
		listenerMock.EXPECT().Accept(gomock.Any()).
			Return(nil, assert.AnError).Times(1),
		listenerMock.EXPECT().Accept(gomock.Any()).
			Return(nil, nil).AnyTimes(),
	)

	var wg sync.WaitGroup
	wg.Add(1)
	called := false
	callbackFn := func(c connection.Connection) {
		assert.NotNil(t, c)
		if !called {
			wg.Done()
			called = true
		}
	}

	l := New(callbackFn).(*listener)
	l.startListenerFn = startListenerFn

	require.NoError(t, l.Start(1111, &tls.Config{}))
	require.EqualError(t, l.Start(1111, &tls.Config{}), ErrAlreadyStarted.Error())

	// wait for the callback to be called and shutdown the listener
	wg.Wait()
	listenerMock.EXPECT().Close().Return(assert.AnError).Times(1)
	require.EqualError(t, l.Shutdown(), assert.AnError.Error())

	for {
		if !l.started {
			break
		}
	}
	require.NoError(t, l.Shutdown())

	// make sure the callback was called
	require.True(t, called)
}
