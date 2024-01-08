package testutil

import (
	"assignment/lib/entity"
	"sync"
)

// MessageCollector is a helper struct for collecting messages.
type MessageCollector struct {
	sync.RWMutex

	messages []entity.Message
}

func (m *MessageCollector) Add(msg entity.Message) {
	m.Lock()
	defer m.Unlock()

	m.messages = append(m.messages, msg)
}

func (m *MessageCollector) Get() []entity.Message {
	m.RLock()
	defer m.RUnlock()

	return m.messages
}

func NewMessageCollector() *MessageCollector {
	return &MessageCollector{
		messages: make([]entity.Message, 0),
	}
}
