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

func TestFlatMap(t *testing.T) {
	tests := []struct {
		name     string
		recv     streams.Connectable
		in       []string
		expected []string
	}{
		{
			name:     "string to string",
			in:       []string{"a", "b", "c"},
			expected: []string{"a", "A", "b", "B", "c", "C"},
			recv: streams.NewFlatMap(func(in string) []string {
				upper := strings.ToUpper(in)
				return []string{in, upper}
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make(chan any, 3)
			out := make(chan any, 9)

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
