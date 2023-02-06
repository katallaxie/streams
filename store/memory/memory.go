package memory

import (
	"sort"
	"sync"

	"github.com/ionos-cloud/streams/store"
)

type memory struct {
	keys    []string
	storage map[string][]byte

	sync.RWMutex
}

// New ...
func New() store.Storage {
	m := new(memory)
	m.keys = make([]string, 0)
	m.storage = make(map[string][]byte)

	return m
}

// Has ...
func (m *memory) Has(key string) (bool, error) {
	m.RLock()
	defer m.RUnlock()

	_, has := m.storage[key]

	return has, nil
}

// Get ...
func (m *memory) Get(key string) ([]byte, error) {
	m.RLock()
	defer m.RUnlock()

	v, ok := m.storage[key]
	if !ok {
		return nil, store.ErrNotExists
	}

	return v, nil
}

// Set ...
func (m *memory) Set(key string, value []byte) error {
	m.Lock()
	defer m.Unlock()

	if value == nil {
		return store.ErrNoNilValue
	}

	if _, ok := m.storage[key]; !ok {
		m.keys = append(m.keys, key)
		sort.Strings(m.keys)
	}

	m.storage[key] = value

	return nil
}

// Delete ...
func (m *memory) Delete(key string) error {
	m.Lock()
	defer m.Unlock()

	delete(m.storage, key)
	for i, k := range m.keys {
		if k == key {
			m.keys = append(m.keys[:i], m.keys[i+1:]...)
			break
		}
	}

	return nil
}

// Open ...
func (m *memory) Open() error {
	return nil
}

// Close ...
func (m *memory) Close() error {
	return nil
}
