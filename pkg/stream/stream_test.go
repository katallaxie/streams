package stream_test

import (
	"testing"

	"github.com/ionos-cloud/streams/pkg/stream"
	"github.com/stretchr/testify/assert"
)

func TestStreamMap(t *testing.T) {
	in := make(chan *stream.Message)

	s := stream.NewStream(in)
	assert.NotNil(t, s)
}
