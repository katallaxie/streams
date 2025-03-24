package main

import (
	"io"
	"strings"

	"github.com/katallaxie/pkg/conv"
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
	if err != nil {
		panic(err)
	}

	sink := sinks.NewStdout()
	s.Pipe(streams.NewPassThrough()).Pipe(streams.NewMap(conv.String)).To(sink)
}
