package streams

// NextCursor ...
type NextCursor struct {
	Key    string
	Value  []byte
	Latest bool
}

// Iterator ...
type Iterator interface {
	// Next ...
	Next() <-chan NextCursor
}

// Table ...
type Table interface {
	// Set ...
	Set(key string, value []byte) error

	// Delete ...
	Delete(key string) error

	// Setup ...
	Setup() error

	// Error ...
	Error() error

	// Sink is the interface that wraps the basic Sink method.
	Sink[string, []byte]

	Iterator
}
