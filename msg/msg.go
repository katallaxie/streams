package msg

import (
	"sync"
)

// Message represents a message in a Stream
type Message[K, V any] interface {
	// Key is used to get the key of a message.
	Key() K

	// Value is used to get the value of a message.
	Mark()

	// Marked is used to check if a message has been marked as processed
	Marked() bool

	// Offset is used to get the offset of a message.
	Offset() int

	// Partition is used to get the partition of a message.
	Partition() int

	// SetKey is used to set the key of a message.
	SetKey(key K)

	// SetValue is used to set the value of a message.
	SetValue(val V)

	// Topic is used to get the topic of a message.
	Topic() string

	// SetTopic is used to set the topic of a message.
	SetTopic(topic string)

	// Value is used to get the value of a message.
	Value() V
}

// Marker is used to mark a message as processed.
type Marker[K, V any] chan Message[K, V]

// MessageImpl is the default implementation of Message
type MessageImpl[K, V any] struct {
	key       K
	marked    bool
	offset    int
	partition int
	topic     string
	val       V

	markedOnce sync.Once

	sync.Mutex
}

// NewMessage creates a new Message.
func NewMessage[K, V any](key K, val V, offset int, partition int, topic string) Message[K, V] {
	return &MessageImpl[K, V]{
		key:       key,
		offset:    offset,
		partition: partition,
		topic:     topic,
		val:       val,
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

// SetValue is used to set the value of a message.
func (m *MessageImpl[K, V]) SetValue(val V) {
	m.val = val
}

// Mark is used to mark a message as processed
func (m *MessageImpl[K, V]) Mark() {
	m.markedOnce.Do(func() {
		m.marked = true
	})
}

// SetTopic is used to set the topic of a message.
func (m *MessageImpl[K, V]) SetTopic(topic string) {
	m.topic = topic
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
