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

// StorageUnimplemented ...
type StorageUnimplemented struct{}

// Open ...
func (s *StorageUnimplemented) Open() error {
	return ErrUnimplemented
}

// Close ...
func (s *StorageUnimplemented) Close() error {
	return ErrUnimplemented
}

// Has ...
func (s *StorageUnimplemented) Has(key string) (bool, error) {
	return false, ErrUnimplemented
}

// Get ...
func (s *StorageUnimplemented) Get(key string) ([]byte, error) {
	return nil, ErrUnimplemented
}

// Set ...
func (s *StorageUnimplemented) Set(key string, value []byte) error {
	return ErrUnimplemented
}

// Delete ...
func (s *StorageUnimplemented) Delete(key string) error {
	return ErrUnimplemented
}
