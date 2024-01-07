package controller

import (
	"assignment/lib/entity"

	"go.uber.org/zap"
)

// notifier is a wrapper around sender (publisher or subscriber)
// with a separate message queue that handles communication to
// that receiver independently.
type notifier struct {
	messages chan entity.Message
	close    chan struct{}
	logger   *zap.Logger
	sender   sender
}

type sender interface {
	SendMessage(entity.Message) error
}

// newNotifier creates a new notifier with the given sender.
func newNotifier(sender sender, logger *zap.Logger) *notifier {
	n := &notifier{
		messages: make(chan entity.Message, DefaultMessageBufferSize),
		close:    make(chan struct{}),
		logger:   logger,
		sender:   sender,
	}

	go n.run()
	return n
}

func (n *notifier) queueMessage(message entity.Message) {
	select {
	case n.messages <- message:
		return
	default:
		n.logger.Warn("Message queue is full, message dropped", zap.String("message", message.Text))
	}
}

func (n *notifier) run() {
	for {
		select {
		case <-n.close:
			return
		case message := <-n.messages:
			if err := n.sender.SendMessage(message); err != nil {
				n.logger.Error("Failed to send message", zap.Error(err), zap.String("message", message.Text))
			}
		}
	}
}

func (n *notifier) stop() {
	n.close <- struct{}{}
}
