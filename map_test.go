package streams_test

import (
	"strings"
	"testing"

	"github.com/katallaxie/streams"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name     string
		recv     streams.Connectable
		in       []string
		expected []string
	}{
		{
			name: "string to string",
			recv: streams.NewMap(strings.ToUpper),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
