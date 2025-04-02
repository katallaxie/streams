package sinks_test

import (
	"testing"

	"github.com/katallaxie/streams"
	"github.com/katallaxie/streams/sinks"
	"github.com/katallaxie/streams/sources"
	"github.com/stretchr/testify/assert"

	"github.com/katallaxie/pkg/fsmx"
)

const (
	AnyAction = iota
)

type NoopState struct {
	Name string
}

func TestStore(t *testing.T) {
	t.Parallel()

	state := NoopState{}

	noopReducer := func(prev fsmx.State, action fsmx.Action) fsmx.State {
		if action.Type() == AnyAction {
			return NoopState{Name: action.Payload().(string)}
		}

		return prev
	}

	s := fsmx.New(state, noopReducer)

	actionMap := func(x any) fsmx.Action {
		return fsmx.NewAction(AnyAction, x.(string))
	}

	in := make(chan any, 1)
	source := sources.NewChanSource(in)
	sink := sinks.NewStore(s, actionMap)

	in <- "foo"
	close(in)

	source.Pipe(streams.NewPassThrough()).To(sink)

	finalState := s.State().(NoopState)
	assert.Equal(t, "foo", finalState.Name)
}
