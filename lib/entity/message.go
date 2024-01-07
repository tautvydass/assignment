package entity

// Message the message abstraction for communication
// between server and clients.
type Message struct {
	Text string
}

// Bytes converts the message to a byte slice.
func (m *Message) Bytes() []byte {
	return []byte(m.Text)
}

// MessageFromBytes converts a byte slice to a message.
func MessageFromBytes(b []byte) Message {
	return Message{
		Text: string(b),
	}
}
