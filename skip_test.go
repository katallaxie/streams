package streams_test

import (
	"testing"

	"github.com/katallaxie/pkg/channels"
	"github.com/katallaxie/streams"
	"github.com/katallaxie/streams/sinks"
	"github.com/katallaxie/streams/sources"
	"github.com/stretchr/testify/require"
)

func TestSkip(t *testing.T) {
	tests := []struct {
		name     string
		recv     streams.Operatable
		in       []int
		expected []int
	}{
		{
			name:     "skip 2",
			in:       []int{1, 2, 3, 4, 5},
			expected: []int{3, 4, 5},
			recv:     streams.NewSkip(2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make(chan any, 5)
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
