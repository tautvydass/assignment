package controller

import (
	"fmt"
	"sync"

	"assignment/lib/connection"
	"assignment/lib/entity"
	"assignment/lib/log"

	"go.uber.org/multierr"
)

const (
	// DefaultMessageBufferSize is the default size of the message buffer.
	DefaultMessageBufferSize = 100
	// MessageNoSubscribers is the message sent to publishers when there are no subscribers.
	MessageNoSubscribers = "No subscribers are currently connected"
	// MessageNewSubscriber is the message sent to publishers when a new subscriber connects.
	MessageNewSubscriber = "New subscriber connected"
	// MessageHelloSubscriber  is the message sent to subscribers when they connect.
	MessageHelloSubscriber = "Hello from server! You're all set."
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
}

// NewCommsController creates a new comms controller.
func NewCommsController() CommsController {
	c := &commsController{
		publishers:  make(map[connection.ReadWriteStream]*notifier),
		subscribers: make(map[connection.WriteStream]*notifier),
		messages:    make(chan entity.Message, DefaultMessageBufferSize),
		close:       make(chan struct{}),
	}

	go c.run()
	return c
}

func (c *commsController) AddPublisher(publisher connection.ReadWriteStream) {
	notifier := newNotifier(publisher, nil)
	publisher.SetConnClosedCallback(func() { c.removePublisher(publisher) })

	c.Lock()
	c.publishers[publisher] = notifier
	c.Unlock()
	log.Info("New publisher successfully connected")

	// Inform the publisher of the current subscriber count.
	message := MessageNoSubscribers
	if subscriberCount := c.subscriberCount(); subscriberCount > 0 {
		message = fmt.Sprintf("%d subscriber(s) currently connected", subscriberCount)
	}
	notifier.queueMessage(entity.Message{Text: message})
}

func (c *commsController) AddSubscriber(subscriber connection.WriteStream) {
	notifier := newNotifier(subscriber, c.removeSubscriber)

	// Say hello to the subscriber to establish the connection.
	// TODO: remove this once WriteStream supports pinging the peer.
	message := entity.Message{Text: MessageHelloSubscriber}
	notifier.queueMessage(message)

	c.Lock()
	c.subscribers[subscriber] = notifier
	c.Unlock()
	log.Info("New subscriber successfully connected")

	// Inform the publishers of the new subscriber.
	message = entity.Message{Text: MessageNewSubscriber}
	c.sendToPublishers(message)
}

func (c *commsController) MessageReceiver() connection.MessageReceiver {
	return func(message entity.Message) {
		select {
		case c.messages <- message:
			return
		default:
			// Too many incoming messages, can't handle them all.
			log.Warnf("Message queue is full, message %q dropped", message.Text)
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
			log.Errorf("Error closing subscriber: %s", err.Error())
			merr = multierr.Append(merr, err)
		}
	}

	for publisherStream, notifier := range c.publishers {
		notifier.stop()
		if err := publisherStream.CloseStream(); err != nil {
			log.Errorf("Error closing publisher: %s", err.Error())
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
			log.Infof("Received message from publisher: %q", msg.Text)
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

func (c *commsController) removePublisher(publisher connection.ReadWriteStream) {
	if err := publisher.CloseStream(); err != nil {
		log.Errorf("Error closing publisher stream: %s", err.Error())
	}

	c.Lock()
	defer c.Unlock()

	notifier, ok := c.publishers[publisher]
	if !ok {
		// This should never happen.
		log.Warn("Notifier for publisher not found")
	} else {
		notifier.stop()
	}

	delete(c.publishers, publisher)
	log.Warn("Publisher disconnected")
}

func (c *commsController) removeSubscriber(sender sender) {
	subscriber, ok := sender.(connection.WriteStream)
	if !ok {
		log.Error("Failed to cast sender to write stream")
		return
	}

	if err := subscriber.CloseStream(); err != nil {
		log.Errorf("Error closing subscriber stream: %s", err.Error())
	}

	c.Lock()
	delete(c.subscribers, subscriber)
	log.Warn("Subscriber disconnected")

	if len(c.subscribers) == 0 {
		// Inform the publishers that there are no subscribers connected.
		c.Unlock()
		message := entity.Message{Text: MessageNoSubscribers}
		c.sendToPublishers(message)
		return
	}

	c.Unlock()
}
