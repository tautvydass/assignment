package controller

import (
	"assignment/lib/connection"
	"assignment/lib/entity"
	"sync"

	"go.uber.org/multierr"
	"go.uber.org/zap"
)

const (
	// DefaultMessageBufferSize is the default size of the message buffer.
	DefaultMessageBufferSize = 100
	// MessageNoSubscribers is the message sent to publishers when there are no subscribers.
	MessageNoSubscribers = "No subscribers are currently connected"
)

// CommsController is the interface for the comms controller. It is responsible
// for managing the communication between publishers and subscribers.
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
	// TODO: implement publisher adding
}

func (c *commsController) AddSubscriber(subscriber connection.WriteStream) {
	// TODO: implement subscriber adding
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
