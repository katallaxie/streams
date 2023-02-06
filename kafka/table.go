package kafka

import "github.com/ionos-cloud/streams/view"

type table struct{}

// NewTable ...
func NewTable() view.Table {
	t := new(table)

	return t
}

// Set ...
func (t *table) Set(key string, value []byte) error {
	return nil
}

// Delete ...
func (t *table) Delete(key string) error {
	return nil
}

// Next ...
func (t *table) Next() <-chan view.NextCursor {
	return nil
}
