package streams_test

import (
	"testing"

	"github.com/katallaxie/pkg/channels"
	"github.com/katallaxie/streams"
	"github.com/katallaxie/streams/sinks"
	"github.com/katallaxie/streams/sources"
	"github.com/stretchr/testify/require"
)

func TestDo(t *testing.T) {
	tests := []struct {
		name     string
		in       []string
		expected []string
	}{
		{
			name:     "string to string",
			in:       []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make(chan any, 3)
			out := make(chan any, 3)

			channels.Channel(tt.in, in)

			s := make([]string, 0)
			recv := streams.Do(func(x string) error { s = append(s, x); return nil })

			source := sources.NewChanSource(in)
			sink := sinks.NewChanSink(out)

			close(in)

			err := source.Pipe(recv).To(sink)
			require.NoError(t, err)

			require.Equal(t, tt.expected, s)
		})
	}
}
