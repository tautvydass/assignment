package controller

import (
	"assignment/lib/entity"
	"assignment/lib/log"
)

// notifier is a wrapper around sender (publisher or subscriber)
// with a separate message queue that handles communication to
// that receiver independently.
type notifier struct {
	messages chan entity.Message
	close    chan struct{}
	sender   sender
}

type sender interface {
	SendMessage(entity.Message) error
}

// newNotifier creates a new notifier with the given sender.
func newNotifier(sender sender) *notifier {
	n := &notifier{
		messages: make(chan entity.Message, DefaultMessageBufferSize),
		close:    make(chan struct{}),
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
		log.Warnf("Message queue is full, message %q dropped", message.Text)
	}
}

func (n *notifier) run() {
	for {
		select {
		case <-n.close:
			return
		case message := <-n.messages:
			if err := n.sender.SendMessage(message); err != nil {
				log.Errorf("Failed to send message %q: %s", message.Text, err.Error())
			}
		}
	}
}

func (n *notifier) stop() {
	n.close <- struct{}{}
}
