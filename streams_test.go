package stream

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockSource struct {
	in  chan *Message
	buf []*Message
}

func (m *mockSource) Messages() chan *Message {
	return m.in
}

func (m *mockSource) Commit() error {
	return nil
}

func newMockSource() *mockSource {
	return &mockSource{
		make(chan *Message),
		make([]*Message, 0),
	}
}

type mockSink struct {
	buf []*Message
}

func (m *mockSink) Write(msg *Message) error {
	m.buf = append(m.buf, msg)

	return nil
}

func newMockSink() *mockSink {
	return &mockSink{
		make([]*Message, 0),
	}
}

func TestStreamMap(t *testing.T) {
	src := newMockSource()

	s := NewStream(src, 0)
	assert.NotNil(t, s)

	out := s.Map(func(m *Message) (*Message, error) {
		m.Name = "foobar"

		return m, nil
	})

	go func() {
		src.in <- &Message{Name: "test"}
	}()

	m := <-out.in
	assert.Equal(t, "foobar", m.Name)
}

func TestStreamMapError(t *testing.T) {
	src := newMockSource()

	err := NewStream(src, 0).Map(func(m *Message) (*Message, error) {
		return nil, fmt.Errorf("error")
	})

	assert.Error(t, fmt.Errorf("error"), err)
}

func TestStreamFilter(t *testing.T) {
	src := newMockSource()

	s := NewStream(src, 0)
	assert.NotNil(t, s)

	out := s.Filter(func(m *Message) (bool, error) {
		return true, nil
	})

	go func() {
		src.in <- &Message{Name: "test"}
	}()

	m := <-out.in
	assert.Equal(t, "test", m.Name)
}

func TestStreamFilterError(t *testing.T) {
	src := newMockSource()

	err := NewStream(src, 0).Filter(func(m *Message) (bool, error) {
		return true, fmt.Errorf("error")
	})

	assert.Error(t, fmt.Errorf("error"), err)
}

func TestStreamBranch(t *testing.T) {
	src := newMockSource()

	s := NewStream(src, 0)
	assert.NotNil(t, s)

	outs := s.Branch(func(m *Message) (bool, error) {
		return true, nil
	}, func(m *Message) (bool, error) {
		return false, nil
	})

	go func() {
		src.in <- &Message{Name: "test"}
	}()

	m := <-outs[0].in
	assert.Equal(t, "test", m.Name)
}

func TestStreamError(t *testing.T) {
	src := newMockSource()

	err := NewStream(src, 0).Branch(func(m *Message) (bool, error) {
		return true, fmt.Errorf("error")
	})

	assert.Error(t, fmt.Errorf("error"), err)
}

func TestStreamSink(t *testing.T) {
	src := newMockSource()
	sink := newMockSink()

	s := NewStream(src, 0)
	assert.NotNil(t, s)

	go func() {
		src.in <- &Message{Name: "test"}
		close(src.in)
	}()

	err := s.Sink(sink)
	assert.NoError(t, err)

	assert.Equal(t, len(sink.buf), 1)
	assert.Equal(t, "test", sink.buf[0].Name)
}
