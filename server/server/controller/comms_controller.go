package controller

import (
	"assignment/lib/connection"
	"assignment/lib/entity"

	"go.uber.org/zap"
)

// DefaultMessageBufferSize is the default size of the message buffer.
const DefaultMessageBufferSize = 100

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
	messages chan entity.Message
	close    chan struct{}
	logger   *zap.Logger
}

// NewCommsController creates a new comms controller.
func NewCommsController(logger *zap.Logger) CommsController {
	c := &commsController{
		messages: make(chan entity.Message, DefaultMessageBufferSize),
		close:    make(chan struct{}),
		logger:   logger,
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
	// TODO: implement close
	return nil
}

func (c *commsController) run() {
	for {
		select {
		case <-c.close:
			return
		case msg := <-c.messages:
			c.publish(msg)
		}
	}
}

func (c *commsController) publish(msg entity.Message) {
	// TODO: implement message publishing
}
