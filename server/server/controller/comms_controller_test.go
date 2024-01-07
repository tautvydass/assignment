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

func TestCommsController_AddPublisher_and_AddSubscriber(t *testing.T) {
	t.Run("subscriber_added_after_publisher", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(2)

		ctrl := gomock.NewController(t)
		publisherStream := connectionmock.NewMockReadWriteStream(ctrl)
		gomock.InOrder(
			publisherStream.EXPECT().SendMessage(entity.Message{
				Text: MessageNoSubscribers,
			}).DoAndReturn(func(_ entity.Message) error {
				wg.Done()
				return nil
			}).Times(1),
			publisherStream.EXPECT().SendMessage(entity.Message{
				Text: MessageNewSubscriber,
			}).DoAndReturn(func(_ entity.Message) error {
				wg.Done()
				return nil
			}).Times(1),
			publisherStream.EXPECT().CloseStream().Return(nil).Times(1),
		)

		subscriberStream := connectionmock.NewMockReadWriteStream(ctrl)
		subscriberStream.EXPECT().CloseStream().Return(nil).Times(1)

		c := NewCommsController(zap.NewNop()).(*commsController)
		defer c.Close()

		c.AddPublisher(publisherStream)
		c.AddSubscriber(subscriberStream)

		wg.Wait()

		require.Len(t, c.publishers, 1)
		require.Len(t, c.subscribers, 1)
	})
	t.Run("publisher_joins_in_between_subscribers", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(2)

		ctrl := gomock.NewController(t)
		publisherStream := connectionmock.NewMockReadWriteStream(ctrl)
		gomock.InOrder(
			publisherStream.EXPECT().SendMessage(entity.Message{
				Text: "1 subscriber(s) currently connected",
			}).DoAndReturn(func(_ entity.Message) error {
				wg.Done()
				return nil
			}).Times(1),
			publisherStream.EXPECT().SendMessage(entity.Message{
				Text: MessageNewSubscriber,
			}).DoAndReturn(func(_ entity.Message) error {
				wg.Done()
				return nil
			}).Times(1),
			publisherStream.EXPECT().CloseStream().Return(nil).Times(1),
		)

		subscriberStream1 := connectionmock.NewMockReadWriteStream(ctrl)
		subscriberStream1.EXPECT().CloseStream().Return(nil).Times(1)

		subscriberStream2 := connectionmock.NewMockReadWriteStream(ctrl)
		subscriberStream2.EXPECT().CloseStream().Return(nil).Times(1)

		c := NewCommsController(zap.NewNop()).(*commsController)
		defer c.Close()

		c.AddSubscriber(subscriberStream1)
		c.AddPublisher(publisherStream)
		c.AddSubscriber(subscriberStream2)

		wg.Wait()

		require.Len(t, c.publishers, 1)
		require.Len(t, c.subscribers, 2)
	})
}
