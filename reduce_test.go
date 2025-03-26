package streams_test

import (
	"testing"

	"github.com/katallaxie/pkg/channels"
	"github.com/katallaxie/streams"
	"github.com/katallaxie/streams/sinks"
	"github.com/katallaxie/streams/sources"
	"github.com/stretchr/testify/require"
)

func sum(a, b int) int {
	return a + b
}

func TestReduce(t *testing.T) {
	tests := []struct {
		name     string
		recv     streams.Operatable
		in       []int
		expected []int
	}{
		{
			name:     "sum",
			in:       []int{1, 2, 3},
			expected: []int{1, 3, 6},
			recv:     streams.NewReduce(sum),
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

			output := channels.Slice[int](out)
			require.Equal(t, tt.expected, output)
		})
	}
}
