package store

import "errors"

// Storage ...
type Storage interface {
	// Open ...
	Open() error

	// Close ...
	Close() error

	// Has ...
	Has(key string) (bool, error)

	// Get ...
	Get(key string) ([]byte, error)

	// Set ...
	Set(key string, value []byte) error

	// Delete ...
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
