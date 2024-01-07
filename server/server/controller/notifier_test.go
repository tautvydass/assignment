package controller

import (
	"sync"
	"testing"

	"assignment/lib/entity"

	"github.com/stretchr/testify/assert"
)

func TestNotifier(t *testing.T) {
	var wg sync.WaitGroup
	sender := newTestSender(func() {
		wg.Done()
	}, []error{nil, nil, nil})
	notifier := newNotifier(sender, nil)

	messages := []entity.Message{
		{Text: "message 1"},
		{Text: "message 2"},
		{Text: "message 3"},
	}
	for _, message := range messages {
		wg.Add(1)
		go notifier.queueMessage(message)
	}

	wg.Wait()
	assert.ElementsMatch(t, messages, sender.messages)

	notifier.stop()
}

type testSender struct {
	sync.RWMutex
	messages  []entity.Message
	responses []error

	callback func()
}

func newTestSender(cb func(), responses []error) *testSender {
	return &testSender{
		messages:  make([]entity.Message, 0, len(responses)),
		responses: responses,
		callback:  cb,
	}
}

func (s *testSender) SendMessage(message entity.Message) error {
	s.Lock()
	defer s.Unlock()
	defer s.callback()

	s.messages = append(s.messages, message)
	if len(s.responses) > 0 {
		response := s.responses[0]
		s.responses = s.responses[1:]
		return response
	}
	return nil
}
