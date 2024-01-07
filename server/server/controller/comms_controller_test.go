package controller

import (
	"sync"
	"testing"

	connectionmock "assignment/lib/connection/mocks"
	"assignment/lib/entity"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCommsController_Close(t *testing.T) {
	c := NewCommsController(zap.NewNop())
	require.NoError(t, c.Close())
}

func TestCommsController_MessageReceiver_and_sendToSubscribers(t *testing.T) {
	c := NewCommsController(zap.NewNop()).(*commsController)
	defer c.Close()

	var wg sync.WaitGroup
	sender := newTestSender(func() { wg.Done() }, nil)
	ctrl := gomock.NewController(t)
	for i := 0; i < 3; i++ {
		streamMock := connectionmock.NewMockReadWriteStream(ctrl)
		streamMock.EXPECT().CloseStream().Return(nil).Times(1)
		c.subscribers[streamMock] = newNotifier(sender, zap.NewNop())
		wg.Add(1)
	}
	require.Len(t, c.subscribers, 3)

	message := entity.Message{Text: "message"}
	go c.MessageReceiver()(message)

	wg.Wait()
	// Make sure that each notifier has received and sent the message.
	require.Equal(t, []entity.Message{message, message, message}, sender.messages)
}
