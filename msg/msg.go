package msg

import "sync"

// Message represents a message in a Stream
type Message interface {
	Mark()
	Marked() bool
	Key() string
	SetKey(key string)
}

// MessageImpl is the default implementation of Message
type MessageImpl struct {
	key        string
	marked     bool
	markedOnce sync.Once

	sync.Mutex
}

// NewMessage creates a new Message.
func NewMessage(key string) Message {
	return &MessageImpl{
		key: key,
	}
}

// Key is used to get the key of a message.
func (m *MessageImpl) Key() string {
	return m.key
}

// SetKey is used to set the key of a message.
func (m *MessageImpl) SetKey(key string) {
	m.key = key
}

// Mark is used to mark a message as processed
func (m *MessageImpl) Mark() {
	m.markedOnce.Do(func() {
		m.marked = true
	})
}

// Marked is used to check if a message has been marked as processed
func (m *MessageImpl) Marked() bool {
	m.Lock()
	defer m.Unlock()

	return m.marked
}
