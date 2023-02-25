package sink

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSink(t *testing.T) {
	tests := []struct {
		desc string
	}{
		{
			desc: "create new sink",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			s := NewSink[string, string]()
			assert.NotNil(t, s)
		})
	}
}

func TestWithContext(t *testing.T) {
	tests := []struct {
		desc string
	}{
		{
			desc: "create new sink with context",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			s := WithContext[string, string](context.Background(), nil, nil, nil)
			assert.NotNil(t, s)
		})
	}
}
