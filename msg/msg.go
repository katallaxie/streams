package msg

import (
	"sync"
)

// Message represents a message in a Stream
type Message[K, V any] interface {
	Mark()
	Marked() bool
	Key() K
	Value() V
	SetKey(key K)
	Topic() string
	Offset() int
	Partition() int
}

// MessageImpl is the default implementation of Message
type MessageImpl[K, V any] struct {
	key        K
	val        V
	partition  int
	offset     int
	topic      string
	marked     bool
	markedOnce sync.Once

	sync.Mutex
}

// NewMessage creates a new Message.
func NewMessage[K, V any](key K, val V, offset int, partition int, topic string) Message[K, V] {
	return &MessageImpl[K, V]{
		key:       key,
		val:       val,
		offset:    offset,
		partition: partition,
		topic:     topic,
	}
}

// Key is used to get the key of a message.
func (m *MessageImpl[K, V]) Key() K {
	return m.key
}

// Value is used to get the value of a message.
func (m *MessageImpl[K, V]) Value() V {
	return m.val
}

// SetKey is used to set the key of a message.
func (m *MessageImpl[K, V]) SetKey(key K) {
	m.key = key
}

// Mark is used to mark a message as processed
func (m *MessageImpl[K, V]) Mark() {
	m.markedOnce.Do(func() {
		m.marked = true
	})
}

// Marked is used to check if a message has been marked as processed
func (m *MessageImpl[K, V]) Marked() bool {
	m.Lock()
	defer m.Unlock()

	return m.marked
}

// Topic ...
func (m *MessageImpl[K, V]) Topic() string {
	return m.topic
}

// Offset ...
func (m *MessageImpl[K, V]) Offset() int {
	return m.offset
}

// Partition ...
func (m *MessageImpl[K, V]) Partition() int {
	return m.partition
}
