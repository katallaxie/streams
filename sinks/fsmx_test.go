package sinks_test

import (
	"testing"

	"github.com/katallaxie/pkg/channels"
	"github.com/katallaxie/pkg/fsmx"
	"github.com/katallaxie/streams"
	"github.com/katallaxie/streams/sinks"
	"github.com/katallaxie/streams/sources"
	"github.com/stretchr/testify/assert"
)

func TestNewFSMStore(t *testing.T) {
	t.Parallel()

	type noopState struct {
		Name string
	}
	state := noopState{}

	reducer := func(prev noopState, action fsmx.Action) noopState {
		if action.Type() == fsmx.ActionType(0) {
			prev.Name = action.Payload().(string)
		}

		return prev
	}

	store := fsmx.New(state, reducer)
	fn := func(payload any) fsmx.Action {
		return fsmx.NewAction(fsmx.ActionType(0), payload)
	}

	in := make(chan any, 1)
	channels.Channel([]string{"foo"}, in)
	source := sources.NewChanSource(in)
	sink := sinks.NewFSMStore(store, fn)

	close(in)

	source.Pipe(streams.NewPassThrough()).To(sink)

	assert.Equal(t, "foo", store.State().Name)
}
