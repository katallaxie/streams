package stream

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamMap(t *testing.T) {
	in := make(chan *Message)

	s := NewStream(in)
	assert.NotNil(t, s)

	out := s.Map(func(m *Message) (*Message, error) {
		m.Name = "foobar"

		return m, nil
	})

	go func() {
		in <- &Message{Name: "test"}
	}()

	m := <-out.in
	assert.Equal(t, "foobar", m.Name)
}

func TestStreamMapError(t *testing.T) {
	in := make(chan *Message)

	err := NewStream(in).Map(func(m *Message) (*Message, error) {
		return nil, fmt.Errorf("error")
	})

	assert.Error(t, fmt.Errorf("error"), err)
}

func TestStreamFilter(t *testing.T) {
	in := make(chan *Message)

	s := NewStream(in)
	assert.NotNil(t, s)

	out := s.Filter(func(m *Message) (bool, error) {
		return true, nil
	})

	go func() {
		in <- &Message{Name: "test"}
	}()

	m := <-out.in
	assert.Equal(t, "test", m.Name)
}

func TestStreamFilterError(t *testing.T) {
	in := make(chan *Message)

	err := NewStream(in).Filter(func(m *Message) (bool, error) {
		return true, fmt.Errorf("error")
	})

	assert.Error(t, fmt.Errorf("error"), err)
}

func TestStreamBranch(t *testing.T) {
	in := make(chan *Message)

	s := NewStream(in)
	assert.NotNil(t, s)

	outs := s.Branch(func(m *Message) (bool, error) {
		return true, nil
	}, func(m *Message) (bool, error) {
		return false, nil
	})

	go func() {
		in <- &Message{Name: "test"}
	}()

	m := <-outs[0].in
	assert.Equal(t, "test", m.Name)
}

func TestStreamError(t *testing.T) {
	in := make(chan *Message)

	err := NewStream(in).Branch(func(m *Message) (bool, error) {
		return true, fmt.Errorf("error")
	})

	assert.Error(t, fmt.Errorf("error"), err)
}
