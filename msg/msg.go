package msg

import "sync"

// Message ...
type Message struct {
	Name       string
	marked     bool
	markedOnce sync.Once

	sync.Mutex
}

// Mark ...
func (m *Message) Mark() {
	m.markedOnce.Do(func() {
		m.marked = true
	})
}

// Marked ...
func (m *Message) Marked() bool {
	m.Lock()
	defer m.Unlock()

	return m.marked
}
