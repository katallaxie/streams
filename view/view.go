package view

import (
	"context"
	"time"

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
type View interface{}

type view struct {
	store store.Storage
	table Table

	server.Listener
}

// New ..
func New[K, V any](table Table, store store.Storage) View {
	v := new(view)
	v.table = table
	v.store = store

	return v
}

// Start ...
func (v *view) Start(ctx context.Context, ready server.ReadyFunc, run server.RunFunc) func() error {
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
