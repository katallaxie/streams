package stream

import (
	"fmt"
	"testing"

	"github.com/ionos-cloud/streams/msg"
	"github.com/stretchr/testify/assert"
)

type mockSource struct {
	in  chan *msg.Message
	buf []*msg.Message
}

func (m *mockSource) Messages() chan *msg.Message {
	return m.in
}

func (m *mockSource) Commit(msgs ...*msg.Message) error {
	m.buf = append(m.buf, msgs...)

	return nil
}

func newMockSource() *mockSource {
	return &mockSource{
		make(chan *msg.Message),
		make([]*msg.Message, 0),
	}
}

type mockSink struct {
	buf []*msg.Message
}

func (m *mockSink) Write(msg ...*msg.Message) error {
	m.buf = append(m.buf, msg...)

	return nil
}

func newMockSink() *mockSink {
	return &mockSink{
		make([]*msg.Message, 0),
	}
}

func TestStreamMap(t *testing.T) {
	src := newMockSource()

	s := NewStream(src, 0)
	assert.NotNil(t, s)

	out := s.Map(func(m *msg.Message) (*msg.Message, error) {
		m.Name = "foobar"

		return m, nil
	})

	go func() {
		src.in <- &msg.Message{Name: "test"}
	}()

	m := <-out.in
	assert.Equal(t, "foobar", m.Name)
}

func TestStreamMapError(t *testing.T) {
	src := newMockSource()

	err := NewStream(src, 0).Map(func(m *msg.Message) (*msg.Message, error) {
		return nil, fmt.Errorf("error")
	})

	assert.Error(t, fmt.Errorf("error"), err)
}

func TestStreamFilter(t *testing.T) {
	src := newMockSource()

	s := NewStream(src, 0)
	assert.NotNil(t, s)

	out := s.Filter(func(m *msg.Message) (bool, error) {
		return true, nil
	})

	go func() {
		src.in <- &msg.Message{Name: "test"}
	}()

	m := <-out.in
	assert.Equal(t, "test", m.Name)
}

func TestStreamFilterError(t *testing.T) {
	src := newMockSource()

	err := NewStream(src, 0).Filter(func(m *msg.Message) (bool, error) {
		return true, fmt.Errorf("error")
	})

	assert.Error(t, fmt.Errorf("error"), err)
}

func TestStreamBranch(t *testing.T) {
	src := newMockSource()

	s := NewStream(src, 0)
	assert.NotNil(t, s)

	outs := s.Branch(func(m *msg.Message) (bool, error) {
		return true, nil
	}, func(m *msg.Message) (bool, error) {
		return false, nil
	})

	go func() {
		src.in <- &msg.Message{Name: "test"}
	}()

	m := <-outs[0].in
	assert.Equal(t, "test", m.Name)
}

func TestStreamError(t *testing.T) {
	src := newMockSource()

	err := NewStream(src, 0).Branch(func(m *msg.Message) (bool, error) {
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
		src.in <- &msg.Message{Name: "test"}
		src.in <- &msg.Message{Name: "test2"}
		close(src.in)
	}()

	err := s.Sink(sink)
	assert.NoError(t, err)

	assert.Equal(t, len(sink.buf), 2)
	assert.Equal(t, "test", sink.buf[0].Name)
	assert.Equal(t, "test2", sink.buf[1].Name)
}

func TestStreamFanOut(t *testing.T) {
	src := newMockSource()

	s := NewStream(src, 0)
	assert.NotNil(t, s)

	outs := s.FanOut(2)

	go func() {
		src.in <- &msg.Message{Name: "test"}
		close(src.in)
	}()

	m := <-outs[0].in
	assert.Equal(t, "test", m.Name)

	m = <-outs[1].in
	assert.Equal(t, "test", m.Name)
}
