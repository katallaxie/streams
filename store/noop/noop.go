package noop

import "github.com/katallaxie/streams/store"

type noop struct {
	store.Unimplemented
}

// NewNull returns a new Null storage.
func New() store.Storage {
	n := new(noop)

	return n
}

// Open is opening the storage.
func (n *noop) Open() error {
	return nil
}

// Close is closing the storage.
func (n *noop) Close() error {
	return nil
}

// Has is checking if a key exists.
func (n *noop) Has(_ string) (bool, error) {
	return false, nil
}

// Get is getting a value from the storage.
func (n *noop) Get(_ string) ([]byte, error) {
	return nil, nil
}

// Set is setting a value in the storage.
func (n *noop) Set(_ string, value []byte) error {
	return nil
}

// Delete is deleting a value from the storage.
func (n *noop) Delete(_ string) error {
	return nil
}
