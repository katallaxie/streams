package noop

import "github.com/ionos-cloud/streams/store"

type noop struct {
	store.StorageUnimplemented
}

// NewNull returns a new Null storage.
func New() store.Storage {
	n := new(noop)

	return n
}

// Open ...
func (n *noop) Open() error {
	return nil
}

// Close ...
func (n *noop) Close() error {
	return nil
}

// Has ...
func (n *noop) Has(key string) (bool, error) {
	return false, nil
}

// Get ...
func (n *noop) Get(key string) ([]byte, error) {
	return nil, nil
}

// Set ...
func (n *noop) Set(key string, value []byte) error {
	return nil
}

// Delete ...
func (n *noop) Delete(key string) error {
	return nil
}
