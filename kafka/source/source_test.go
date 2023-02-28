package source

import (
	"context"
	"testing"
	"time"

	"github.com/ionos-cloud/streams/codec"
	kgo "github.com/segmentio/kafka-go"

	"github.com/stretchr/testify/assert"
)

func TestDefaultOpts(t *testing.T) {
	tests := []struct {
		name string
		want *Opts
	}{
		{
			name: "default options",
			want: &Opts{
				commitMode:    CommitManual,
				bufferSize:    1000,
				bufferTimeout: time.Second,
			},
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
		ctx        context.Context
		r          *kgo.Reader
		key        codec.Decoder[string]
		value      codec.Decoder[string]
		keyEncoder codec.Encoder[string]
		opts       []Opt
	}
	tests := []struct {
		name string
		args args
		want *Source[string, string]
	}{
		{
			name: "with context",
			args: args{
				ctx:        context.Background(),
				r:          nil,
				key:        nil,
				value:      nil,
				keyEncoder: nil,
				opts:       nil,
			},
			want: &Source[string, string]{
				ctx:          context.Background(),
				reader:       nil,
				keyDecoder:   nil,
				valueDecoder: nil,
				keyEncoder:   nil,
				opts:         DefaultOpts(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithContext(tt.args.ctx, tt.args.r, tt.args.key, tt.args.value, tt.args.keyEncoder, tt.args.opts...)
			assert.Equal(t, tt.want, got)
		})
	}
}
