package controller

import (
	"fmt"
	"sync"

	"assignment/lib/connection"
	"assignment/lib/entity"

	"go.uber.org/multierr"
	"go.uber.org/zap"
)

const (
	// DefaultMessageBufferSize is the default size of the message buffer.
	DefaultMessageBufferSize = 100
	// MessageNoSubscribers is the message sent to publishers when there are no subscribers.
	MessageNoSubscribers = "No subscribers are currently connected"
	// MessageNewSubscriber is the message sent to publishers when a new subscriber connects.
	MessageNewSubscriber = "New subscriber connected"
)

// CommsController is the interface for the comms controller. It is responsible
// for managing the communication between publishers and subscribers.
// TODO: handle disconnected publishers and subscribers.
type CommsController interface {
	// AddPublisher adds a publisher to the comms controller.
	AddPublisher(publisher connection.ReadWriteStream)
	// AddSubscriber adds a subscriber to the comms controller.
	AddSubscriber(subscriber connection.WriteStream)
	// MessageReceiver returns a message receiver function for new message handling.
	MessageReceiver() connection.MessageReceiver
	// Close closes the comms controller.
	Close() error
}

type commsController struct {
	sync.RWMutex
	publishers  map[connection.ReadWriteStream]*notifier
	subscribers map[connection.WriteStream]*notifier

	messages chan entity.Message
	close    chan struct{}
	logger   *zap.Logger
}

// NewCommsController creates a new comms controller.
func NewCommsController(logger *zap.Logger) CommsController {
	c := &commsController{
		publishers:  make(map[connection.ReadWriteStream]*notifier),
		subscribers: make(map[connection.WriteStream]*notifier),
		messages:    make(chan entity.Message, DefaultMessageBufferSize),
		close:       make(chan struct{}),
		logger:      logger,
	}

	go c.run()
	return c
}

func (c *commsController) AddPublisher(publisher connection.ReadWriteStream) {
	notifier := newNotifier(publisher, c.logger)

	c.Lock()
	c.publishers[publisher] = notifier
	c.Unlock()
	c.logger.Info("Added new publisher")

	// Inform the publisher of the current subscriber count.
	message := MessageNoSubscribers
	if subscriberCount := c.subscriberCount(); subscriberCount > 0 {
		message = fmt.Sprintf("%d subscriber(s) currently connected", subscriberCount)
	}
	notifier.queueMessage(entity.Message{Text: message})
}

func (c *commsController) AddSubscriber(subscriber connection.WriteStream) {
	notifier := newNotifier(subscriber, c.logger)

	c.Lock()
	c.subscribers[subscriber] = notifier
	c.Unlock()
	c.logger.Info("Added new subscriber")

	// Inform the publishers of the new subscriber.
	message := entity.Message{Text: MessageNewSubscriber}
	c.sendToPublishers(message)
}

func (c *commsController) MessageReceiver() connection.MessageReceiver {
	return func(message entity.Message) {
		select {
		case c.messages <- message:
			return
		default:
			c.logger.Warn("Message queue is full, message dropped", zap.String("message", message.Text))
		}
	}
}

func (c *commsController) Close() error {
	c.close <- struct{}{}
	c.Lock()
	defer c.Unlock()

	var merr error
	for subscriberStream, notifier := range c.subscribers {
		notifier.stop()
		if err := subscriberStream.CloseStream(); err != nil {
			c.logger.Error("Error closing subscriber", zap.Error(err))
			merr = multierr.Append(merr, err)
		}
	}

	for publisherStream, notifier := range c.publishers {
		notifier.stop()
		if err := publisherStream.CloseStream(); err != nil {
			c.logger.Error("Error closing publisher", zap.Error(err))
			merr = multierr.Append(merr, err)
		}
	}

	return merr
}

func (c *commsController) run() {
	for {
		select {
		case <-c.close:
			return
		case msg := <-c.messages:
			c.logger.Info("Received message", zap.String("message", msg.Text))
			c.sendToSubscribers(msg)
		}
	}
}

func (c *commsController) sendToSubscribers(msg entity.Message) {
	notifiers := c.getSubscriberNotifiers()
	for _, notifier := range notifiers {
		notifier.queueMessage(msg)
	}
}

func (c *commsController) getSubscriberNotifiers() []*notifier {
	c.RLock()
	defer c.RUnlock()

	notifiers := make([]*notifier, 0, len(c.subscribers))
	for _, notifier := range c.subscribers {
		notifiers = append(notifiers, notifier)
	}

	return notifiers
}

func (c *commsController) subscriberCount() int {
	c.RLock()
	defer c.RUnlock()

	return len(c.subscribers)
}

func (c *commsController) sendToPublishers(msg entity.Message) {
	notifiers := c.getPublisherNotifiers()
	for _, notifier := range notifiers {
		notifier.queueMessage(msg)
	}
}

func (c *commsController) getPublisherNotifiers() []*notifier {
	c.RLock()
	defer c.RUnlock()

	notifiers := make([]*notifier, 0, len(c.publishers))
	for _, notifier := range c.publishers {
		notifiers = append(notifiers, notifier)
	}

	return notifiers
}
