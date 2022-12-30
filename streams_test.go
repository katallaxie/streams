package streams

import (
	"fmt"
	"testing"

	"github.com/ionos-cloud/streams/msg"
	"github.com/stretchr/testify/assert"
)

type mockSource[K, V any] struct {
	in  chan msg.Message[K, V]
	buf []msg.Message[K, V]
}

func (m *mockSource[K, V]) Messages() chan msg.Message[K, V] {
	return m.in
}

func (m *mockSource[K, V]) Commit(msgs ...msg.Message[K, V]) error {
	m.buf = append(m.buf, msgs...)

	return nil
}

func newMockSource[K, V any]() *mockSource[K, V] {
	return &mockSource[K, V]{
		make(chan msg.Message[K, V]),
		make([]msg.Message[K, V], 0),
	}
}

type mockSink[K, V any] struct {
	buf []msg.Message[K, V]
}

func (m *mockSink[K, V]) Write(msg ...msg.Message[K, V]) error {
	m.buf = append(m.buf, msg...)

	return nil
}

func newMockSink[K, V any]() *mockSink[K, V] {
	return &mockSink[K, V]{
		make([]msg.Message[K, V], 0),
	}
}

func TestStreamMap(t *testing.T) {
	src := newMockSource[string, string]()

	s := NewStream[string, string](src)
	assert.NotNil(t, s)

	out := s.Map(func(m msg.Message[string, string]) (msg.Message[string, string], error) {
		m.SetKey("foobar")

		return m, nil
	})

	go func() {
		src.in <- msg.NewMessage("test", "test")
	}()

	m := <-out.in
	assert.Equal(t, "foobar", m.Key())
}

func TestStreamMapError(t *testing.T) {
	src := newMockSource[string, string]()

	err := NewStream[string, string](src).Map(func(m msg.Message[string, string]) (msg.Message[string, string], error) {
		return nil, fmt.Errorf("error")
	})

	assert.Error(t, fmt.Errorf("error"), err)
}

func TestStreamFilter(t *testing.T) {
	src := newMockSource[string, string]()

	s := NewStream[string, string](src)
	assert.NotNil(t, s)

	out := s.Filter(func(m msg.Message[string, string]) (bool, error) {
		return true, nil
	})

	go func() {
		src.in <- msg.NewMessage("test", "test")
	}()

	m := <-out.in
	assert.Equal(t, "test", m.Key())
}

func TestStreamFilterError(t *testing.T) {
	src := newMockSource[string, string]()

	err := NewStream[string, string](src).Filter(func(m msg.Message[string, string]) (bool, error) {
		return true, fmt.Errorf("error")
	})

	assert.Error(t, fmt.Errorf("error"), err)
}

func TestStreamBranch(t *testing.T) {
	src := newMockSource[string, string]()

	s := NewStream[string, string](src)
	assert.NotNil(t, s)

	outs := s.Branch(func(m msg.Message[string, string]) (bool, error) {
		return true, nil
	}, func(m msg.Message[string, string]) (bool, error) {
		return false, nil
	})

	go func() {
		src.in <- msg.NewMessage("test", "test")
	}()

	m := <-outs[0].in
	assert.Equal(t, "test", m.Key())
}

func TestStreamError(t *testing.T) {
	src := newMockSource[string, string]()

	err := NewStream[string, string](src).Branch(func(m msg.Message[string, string]) (bool, error) {
		return true, fmt.Errorf("error")
	})

	assert.Error(t, fmt.Errorf("error"), err)
}

func TestStreamSink(t *testing.T) {
	src := newMockSource[string, string]()
	sink := newMockSink[string, string]()

	s := NewStream[string, string](src)
	assert.NotNil(t, s)

	go func() {
		src.in <- msg.NewMessage("test", "test")
		src.in <- msg.NewMessage("test2", "test")
		close(src.in)
	}()

	err := s.Sink(sink)
	assert.NoError(t, err)

	assert.Equal(t, len(sink.buf), 2)
	assert.Equal(t, "test", sink.buf[0].Key())
	assert.Equal(t, "test2", sink.buf[1].Key())
}

func TestStreamFanOut(t *testing.T) {
	src := newMockSource[string, string]()

	s := NewStream[string, string](src)
	assert.NotNil(t, s)

	outs := s.FanOut(2)

	go func() {
		src.in <- msg.NewMessage("test", "test")
		close(src.in)
	}()

	m := <-outs[0].in
	assert.Equal(t, "test", m.Key())

	m = <-outs[1].in
	assert.Equal(t, "test", m.Key())
}
