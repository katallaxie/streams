package streams_test

import (
	"strings"
	"testing"

	"github.com/katallaxie/pkg/channels"
	"github.com/katallaxie/streams"
	"github.com/katallaxie/streams/sinks"
	"github.com/katallaxie/streams/sources"
	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name     string
		recv     streams.Operatable
		in       []string
		expected []string
	}{
		{
			name:     "string to string",
			in:       []string{"a", "b", "c"},
			expected: []string{"A", "B", "C"},
			recv:     streams.NewMap(strings.ToUpper),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make(chan any, 3)
			out := make(chan any, 3)

			channels.Channel(tt.in, in)

			source := sources.NewChanSource(in)
			sink := sinks.NewChanSink(out)

			close(in)

			source.Pipe(tt.recv).To(sink)

			output := channels.Slice[string](out)
			require.Equal(t, tt.expected, output)
		})
	}
}
