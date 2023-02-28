package streams

import (
	"fmt"
	"sync"
	"testing"

	"github.com/ionos-cloud/streams/msg"

	"github.com/stretchr/testify/assert"
)

type mockSource[K, V any] struct {
	in  chan msg.Message[K, V]
	buf []msg.Message[K, V]

	sync.RWMutex
}

func (m *mockSource[K, V]) Messages() chan msg.Message[K, V] {
	return m.in
}

func (m *mockSource[K, V]) Commit(msgs ...msg.Message[K, V]) error {
	m.Lock()
	defer m.Unlock()

	m.buf = append(m.buf, msgs...)

	return nil
}

func (m *mockSource[K, V]) Error() error {
	return nil
}

func newMockSource[K, V any]() *mockSource[K, V] {
	return &mockSource[K, V]{
		make(chan msg.Message[K, V]),
		make([]msg.Message[K, V], 0),
		sync.RWMutex{},
	}
}

type mockSink[K, V any] struct {
	buf []msg.Message[K, V]
	sync.RWMutex
}

func (m *mockSink[K, V]) Write(msg ...msg.Message[K, V]) error {
	m.Lock()
	defer m.Unlock()

	m.buf = append(m.buf, msg...)

	return nil
}

func newMockSink[K, V any]() *mockSink[K, V] {
	return &mockSink[K, V]{
		make([]msg.Message[K, V], 0),
		sync.RWMutex{},
	}
}

func TestStreamMap(t *testing.T) {
	src := newMockSource[string, string]()

	s := NewStream[string, string](src)
	assert.NotNil(t, s)

	out := s.Map("foo", func(m msg.Message[string, string]) (msg.Message[string, string], error) {
		m.SetKey("foobar")

		return m, nil
	})

	go func() {
		src.in <- msg.NewMessage("test", "test", 0, 0, "", nil)
	}()

	m := <-out.in
	assert.Equal(t, "foobar", m.Key())
}

func TestStreamMapError(t *testing.T) {
	src := newMockSource[string, string]()

	err := NewStream[string, string](src).Map("foo", func(m msg.Message[string, string]) (msg.Message[string, string], error) {
		return nil, fmt.Errorf("error")
	})

	assert.Error(t, fmt.Errorf("error"), err)
}

func TestStreamFilter(t *testing.T) {
	src := newMockSource[string, string]()

	s := NewStream[string, string](src)
	assert.NotNil(t, s)

	out := s.Filter("foo", func(m msg.Message[string, string]) (bool, error) {
		return true, nil
	})

	go func() {
		src.in <- msg.NewMessage("test", "test", 0, 0, "", nil)
	}()

	m := <-out.in
	assert.Equal(t, "test", m.Key())
}

func TestStreamFilterError(t *testing.T) {
	src := newMockSource[string, string]()

	err := NewStream[string, string](src).Filter("foo", func(m msg.Message[string, string]) (bool, error) {
		return true, fmt.Errorf("error")
	})

	assert.Error(t, fmt.Errorf("error"), err)
}

func TestStreamBranch(t *testing.T) {
	src := newMockSource[string, string]()

	s := NewStream[string, string](src)
	assert.NotNil(t, s)

	outs := s.Branch("foo", func(m msg.Message[string, string]) (bool, error) {
		return true, nil
	}, func(m msg.Message[string, string]) (bool, error) {
		return false, nil
	})

	go func() {
		src.in <- msg.NewMessage("test", "test", 0, 0, "", nil)
	}()

	m := <-outs[0].in
	assert.Equal(t, "test", m.Key())
}

func TestStreamError(t *testing.T) {
	src := newMockSource[string, string]()

	err := NewStream[string, string](src).Branch("foo", func(m msg.Message[string, string]) (bool, error) {
		return true, fmt.Errorf("error")
	})

	assert.Error(t, fmt.Errorf("error"), err)
}

func TestStreamSink(t *testing.T) {
	src := newMockSource[string, string]()
	sink := newMockSink[string, string]()

	s := NewStream[string, string](src, WithBuffer(0), WithTimeout(0))
	assert.NotNil(t, s)

	go func() {
		src.in <- msg.NewMessage("test", "test", 0, 0, "", nil)
		src.in <- msg.NewMessage("test2", "test", 0, 0, "", nil)
		close(src.in)
	}()

	s.Sink("foo", sink)
	assert.NoError(t, s.Error())

	assert.Len(t, sink.buf, 2)
	assert.Equal(t, "test", sink.buf[0].Key())
	assert.Equal(t, "test2", sink.buf[1].Key())
}

func TestStreamFanOut(t *testing.T) {
	src := newMockSource[string, string]()

	s := NewStream[string, string](src)
	assert.NotNil(t, s)

	outs := s.FanOut("foo", 2)

	go func() {
		src.in <- msg.NewMessage("test", "test", 0, 0, "", nil)
		close(src.in)
	}()

	m := <-outs[0].in
	assert.Equal(t, "test", m.Key())

	m = <-outs[1].in
	assert.Equal(t, "test", m.Key())
}
