package listener

import (
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
		startListenerFn = func(port int) (QUICListener, error) {
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

	require.NoError(t, l.Start(1111))
	require.EqualError(t, l.Start(1111), ErrAlreadyStarted.Error())

	// wait for the callback to be called and shutdown the listener
	wg.Wait()
	listenerMock.EXPECT().Close().Return(assert.AnError).Times(1)
	require.EqualError(t, l.Shutdown(), assert.AnError.Error())
	require.NoError(t, l.Shutdown())

	// make sure the callback was called
	require.True(t, called)
}
