package store

import "errors"

// Storage ...
type Storage interface {
	// Open is called when the storage is opened.
	Open() error

	// Close is called when the storage is closed.
	Close() error

	// Has is called to check if a key exists.
	Has(key string) (bool, error)

	// Get is called to get a value.
	Get(key string) ([]byte, error)

	// Set is called to set a value.
	Set(key string, value []byte) error

	// Delete is called to delete a value.
	Delete(key string) error
}

var (
	// ErrUnimplemented ...
	ErrUnimplemented = errors.New("not implemented")

	// ErrNoNilValue ...
	ErrNoNilValue = errors.New("no nil value")

	// ErrNotExists ...
	ErrNotExists = errors.New("not exists")
)

// Unimplemented ...
type Unimplemented struct{}

// Open ...
func (s *Unimplemented) Open() error {
	return ErrUnimplemented
}

// Close ...
func (s *Unimplemented) Close() error {
	return ErrUnimplemented
}

// Has ...
func (s *Unimplemented) Has(key string) (bool, error) {
	return false, ErrUnimplemented
}

// Get ...
func (s *Unimplemented) Get(key string) ([]byte, error) {
	return nil, ErrUnimplemented
}

// Set ...
func (s *Unimplemented) Set(key string, value []byte) error {
	return ErrUnimplemented
}

// Delete ...
func (s *Unimplemented) Delete(key string) error {
	return ErrUnimplemented
}
