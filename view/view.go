package view

import "github.com/ionos-cloud/streams/store"

// Table ...
type Table string

// Group ...
type Group string

// GroupTable ...
func GroupTable(group Group) Table {
	return Table(group)
}

// View ...
type View interface{}

type view struct {
	store store.Storage
	table Table
}

// New ..
func New(table Table, store store.Storage) View {
	v := new(view)
	v.table = table
	v.store = store

	return v
}
