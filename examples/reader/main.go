package main

import (
	"io"
	"strings"

	"github.com/katallaxie/pkg/conv"
	"github.com/katallaxie/pkg/errorx"
	"github.com/katallaxie/streams"
	"github.com/katallaxie/streams/sinks"
	"github.com/katallaxie/streams/sources"
)

func read(r io.Reader) ([]byte, error) {
	buf := make([]byte, 1)
	n, err := r.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}

func main() {
	r := io.NopCloser(strings.NewReader("Hello, world!"))

	s, err := sources.NewReaderSource(r, read)
	errorx.Panic(err)

	s.Pipe(streams.DefaultPassThrough).Pipe(streams.NewMap(conv.String)).To(sinks.DefaultStdout)
}
