package sinks_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/katallaxie/streams"
	"github.com/katallaxie/streams/sinks"
	"github.com/katallaxie/streams/sources"

	"github.com/katallaxie/pkg/channels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type NoopWriterCloser struct {
	*bufio.Writer
}

func (w *NoopWriterCloser) Close() error {
	return w.Flush()
}

func TestNewWriter(t *testing.T) {
	t.Parallel()

	in := make(chan any, 1)

	b := bytes.NewBuffer(nil)
	bw := bufio.NewWriter(b)
	mw := &NoopWriterCloser{bw}

	channels.Channel([]string{"foo"}, in)
	source := sources.NewChanSource(in)
	sink, err := sinks.NewWriterSink(mw)
	require.NoError(t, err)

	close(in)

	source.Pipe(streams.NewPassThrough()).To(sink)

	assert.Equal(t, "foo", b.String())
}
