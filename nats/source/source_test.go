package source

import (
	"context"
	"testing"

	"github.com/ionos-cloud/streams/codec"
	"github.com/nats-io/nats.go"

	"github.com/stretchr/testify/assert"
)

func TestDefaultOpts(t *testing.T) {
	tests := []struct {
		name string
		want *Opts
	}{
		{
			name: "default options",
			want: &Opts{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, DefaultOpts())
		})
	}
}

func TestWithContext(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   codec.Decoder[string]
		value codec.Decoder[string]
		sub   *nats.Subscription
		opts  []Opt
	}
	tests := []struct {
		name string
		args args
		want *Source[string, string]
	}{
		{
			name: "with context",
			args: args{
				ctx:   context.Background(),
				sub:   nil,
				key:   nil,
				value: nil,
			},
			want: &Source[string, string]{
				ctx:          context.Background(),
				sub:          nil,
				keyDecoder:   nil,
				valueDecoder: nil,
				opts:         DefaultOpts(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithContext(tt.args.ctx, tt.args.sub, tt.args.key, tt.args.value, tt.args.opts...)
			assert.Equal(t, tt.want, got)
		})
	}
}
