package view

import (
	"context"
	"errors"
	"sync"

	"github.com/katallaxie/streams"
	"github.com/katallaxie/streams/codec"
	"github.com/katallaxie/streams/store"

	"github.com/katallaxie/pkg/server"
	"github.com/katallaxie/pkg/utils"
)

const (
	// Pure ...
	Pure = iota

	// Initializing ...
	Initializing

	// Running ...
	Running
)

// ErrCatchup is returned when the view is not yet caught up with the latest changes.
var ErrCatchup = errors.New("views: catching up with the latest changes")

// Value is a value in the table.
type Value interface {
	int | ~string | []byte
}

// View is a view of the data in the table
type View[V Value] interface {
	// Get is used to retrieve a value from the view.
	Get(key string) (V, error)

	// Set is used to set a value in the view.
	Set(key string, value V) error

	// Delete is used to delete a value from the view.
	Delete(key string) error

	server.Listener
}

type view[V Value] struct {
	store store.Storage
	table streams.Table

	encoder codec.Encoder[V]
	decoder codec.Decoder[V]

	catchUpOnce sync.Once
	catchUp     bool
}

// New initializes a new view.
func New[V Value](table streams.Table, encoder codec.Encoder[V], decoder codec.Decoder[V], store store.Storage) View[V] {
	v := new(view[V])
	v.table = table
	v.store = store

	v.encoder = encoder
	v.decoder = decoder

	return v
}

// Get is used to retrieve a value from the view.
func (v *view[V]) Get(key string) (V, error) {
	if !v.catchUp {
		return utils.Zero[V](), ErrCatchup
	}

	b, err := v.store.Get(key)
	if err != nil {
		return utils.Zero[V](), err
	}

	value, err := v.decoder.Decode(b)
	if err != nil {
		return utils.Zero[V](), err
	}

	return value, nil
}

// Set is used to set a value in the view.
func (v *view[V]) Set(key string, value V) error {
	b, err := v.encoder.Encode(value)
	if err != nil {
		return err
	}

	// TODO: this is non optimistic, the update is published to the table and then synced to storage
	err = v.table.Set(key, b)
	if err != nil {
		return err
	}

	return nil
}

// Delete is used to delete a value from the view.
func (v *view[V]) Delete(key string) error {
	err := v.table.Delete(key) // Tombstone message
	if err != nil {
		return err
	}

	return nil
}

// Start is used to start the view.
func (v *view[V]) Start(ctx context.Context, ready server.ReadyFunc, run server.RunFunc) func() error {
	return func() error {
		err := v.table.Setup()
		if err != nil {
			return err
		}

		ready()

		for c := range v.table.Next() {
			if c.Value == nil {
				ok, err := v.store.Has(c.Key)
				if err != nil {
					return err
				}

				if !ok {
					continue
				}

				err = v.store.Delete(c.Key)
				if err != nil {
					return err
				}

				continue
			}

			err := v.store.Set(c.Key, c.Value)
			if err != nil {
				return err
			}

			if c.Latest {
				v.catchUpOnce.Do(func() {
					v.catchUp = true
				})
			}
		}

		return nil
	}
}
