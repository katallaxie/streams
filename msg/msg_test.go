package msg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	tests := []struct {
		desc      string
		val       string
		key       string
		offset    int
		partition int
		topic     string
	}{
		{
			desc:      "create new message",
			val:       "foo",
			key:       "bar",
			offset:    0,
			partition: 0,
			topic:     "foo",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			m := NewMessage(tc.key, tc.val, tc.offset, tc.partition, tc.topic, nil)
			assert.NotNil(t, m)

			assert.Equal(t, tc.key, m.Key())
			assert.Equal(t, tc.val, m.Value())
			assert.Equal(t, tc.partition, m.Partition())
			assert.Equal(t, tc.offset, m.Offset())
			assert.Equal(t, tc.topic, m.Topic())
		})
	}
}

func TestMarkMessage(t *testing.T) {
	tests := []struct {
		desc      string
		val       string
		key       string
		offset    int
		partition int
		topic     string
	}{
		{
			desc:      "mark message",
			val:       "foo",
			key:       "bar",
			offset:    0,
			partition: 0,
			topic:     "foo",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			mark := make(Marker[string, string], 1)

			m := NewMessage(tc.key, tc.val, tc.offset, tc.partition, tc.topic, mark)
			assert.NotNil(t, m)

			m.Mark()
			x := <-mark

			assert.True(t, m.Marked())
			assert.Equal(t, tc.key, x.Key())
		})
	}
}

func TestSetValue(t *testing.T) {
	tests := []struct {
		desc string
		val  string
	}{
		{
			desc: "set value",
			val:  "foo",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			m := NewMessage("", "", 0, 0, "", nil)
			assert.NotNil(t, m)

			m.SetValue(tc.val)

			assert.Equal(t, tc.val, m.Value())
		})
	}
}

func TestSetKey(t *testing.T) {
	tests := []struct {
		desc string
		key  string
	}{
		{
			desc: "set key",
			key:  "foo",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			m := NewMessage("", "", 0, 0, "", nil)
			assert.NotNil(t, m)

			m.SetKey(tc.key)

			assert.Equal(t, tc.key, m.Key())
		})
	}
}

func TestSetTopic(t *testing.T) {
	tests := []struct {
		desc  string
		topic string
	}{
		{
			desc:  "set topic",
			topic: "foo",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			m := NewMessage("", "", 0, 0, "", nil)
			assert.NotNil(t, m)

			m.SetTopic(tc.topic)

			assert.Equal(t, tc.topic, m.Topic())
		})
	}
}

func TestClone(t *testing.T) {
	tests := []struct {
		desc      string
		val       string
		key       string
		offset    int
		partition int
		topic     string
		mark      Marker[string, string]
	}{
		{
			desc:      "clone message",
			val:       "foo",
			key:       "bar",
			offset:    0,
			partition: 0,
			topic:     "foo",
			mark:      make(Marker[string, string], 1),
		},
		{
			desc:      "clone message with nil marker",
			val:       "foo",
			key:       "bar",
			offset:    0,
			partition: 0,
			topic:     "foo",
			mark:      nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			m := NewMessage(tc.key, tc.val, tc.offset, tc.partition, tc.topic, tc.mark)
			assert.NotNil(t, m)

			c := m.Clone()
			assert.NotNil(t, c)

			assert.Equal(t, tc.key, c.Key())
			assert.Equal(t, tc.val, c.Value())
			assert.Equal(t, tc.partition, c.Partition())
			assert.Equal(t, tc.offset, c.Offset())
			assert.Equal(t, tc.topic, c.Topic())

			if tc.mark != nil {
				m.Mark()
				x := <-tc.mark
				assert.True(t, m.Marked())
				assert.Equal(t, tc.key, x.Key())
			}
		})
	}
}
