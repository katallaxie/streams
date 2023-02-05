package view

import (
	"context"
	"time"

	"github.com/ionos-cloud/streams"
	"github.com/ionos-cloud/streams/store"
	"github.com/katallaxie/pkg/server"
)

const (
	// Pure ...
	Pure = iota

	// Initializing ...
	Initializing

	// Running ...
	Running
)

// NextCursor ...
type NextCursor struct {
	Key   string
	Value []byte
}

// Iterator ...
type Iterator interface {
	// Next ...
	Next() <-chan NextCursor
}

// Table ...
type Table interface {
	Iterator
}

// Unimplemented ...
type Unimplemented struct{}

// Next ...
func (u *Unimplemented) Next() (string, string) {
	return "", ""
}

// View ...
type View[V any] interface {
	// Get ...
	Get(key string) (V, error)

	// Set ...
	Set(key string, value V) error

	// Delete ...
	Delete(key string) error
}

type view[V any] struct {
	store store.Storage
	table Table

	encoder streams.Encoder[V]
	decoder streams.Decoder[V]

	server.Listener
}

// New ..
func New[V any](table Table, encoder streams.Encoder[V], decoder streams.Decoder[V], store store.Storage) View[V] {
	v := new(view[V])
	v.table = table
	v.store = store

	v.encoder = encoder
	v.decoder = decoder

	return v
}

// Get ...
func (v *view[V]) Get(key string) (V, error) {
	b, err := v.store.Get(key)
	if err != nil {
		return Zero[V](), err
	}

	value, err := v.decoder.Decode(b)
	if err != nil {
		return Zero[V](), err
	}

	return value, nil
}

// Set ...
func (v *view[V]) Set(key string, value V) error {
	return nil
}

// Delete ...
func (v *view[V]) Delete(key string) error {
	return nil
}

// Start ...
func (v *view[V]) Start(ctx context.Context, ready server.ReadyFunc, run server.RunFunc) func() error {
	return func() error {
		ready()

		ticker := time.NewTicker(1 * time.Second)

		for {
			select {
			case <-ticker.C:
			case <-ctx.Done():
				return nil
			case c := <-v.table.Next():
				err := v.store.Set(c.Key, c.Value)
				if err != nil {
					return err
				}

				continue
			}
		}
	}
}

// Zero ..
func Zero[T any]() T {
	return *new(T)
}
