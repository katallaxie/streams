package streams

import "errors"

// ErrNotImplemented is returned when a method is not implemented.
var ErrNotImplemented = errors.New("not implemented")

// NextCursor is the next cursor.
type NextCursor struct {
	Key    string
	Value  []byte
	Latest bool
}

// Iterator is the interface that wraps the basic Next method.
type Iterator interface {
	// Next moves the cursor to the next key/value pair, which will then be available through the Key, Value and Latest methods.
	// It returns false if the iterator is exhausted.
	Next() <-chan NextCursor
}

// Table is the interface that wraps the basic Set, Delete, Setup, Error and Sink methods.
type Table interface {
	// Set is setting a key/value pair.
	Set(key string, value []byte) error

	// Delete is deleting a key/value pair.
	Delete(key string) error

	// Setup is setting up the table.
	Setup() error

	// Error is returning the error.
	Error() error

	// Sink is the interface that wraps the basic Sink method.
	Sink[string, []byte]

	Iterator
}

type tableUnimplemented struct{}

// Setup is setting key/value pair.
func (t *tableUnimplemented) Set(key string, value []byte) error {
	return ErrNotImplemented
}

// Delete is deleting a key/value pair.
func (t *tableUnimplemented) Delete(key string) error {
	return ErrNotImplemented
}

// Setup is setting up the table.
func (t *tableUnimplemented) Setup() error {
	return ErrNotImplemented
}

// Error is returning the error.
func (t *tableUnimplemented) Error() error {
	return ErrNotImplemented
}
