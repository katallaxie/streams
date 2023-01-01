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
			m := NewMessage(tc.key, tc.val, tc.offset, tc.partition, tc.topic)
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
			m := NewMessage(tc.key, tc.val, tc.offset, tc.partition, tc.topic)
			assert.NotNil(t, m)

			m.Mark()

			assert.True(t, m.Marked())
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
			m := NewMessage("", "", 0, 0, "")
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
			m := NewMessage("", "", 0, 0, "")
			assert.NotNil(t, m)

			m.SetKey(tc.key)

			assert.Equal(t, tc.key, m.Key())
		})
	}
}
