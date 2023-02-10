package leveldb

import (
	"github.com/ionos-cloud/streams/store"

	"github.com/katallaxie/pkg/utils/files"
	level "github.com/syndtr/goleveldb/leveldb"
)

type leveldb struct {
	db       *level.DB
	teardown files.TearDownFunc

	store.Unimplemented
}

// New ...
func New() store.Storage {
	l := new(leveldb)

	return l
}

// Has ...
func (l *leveldb) Has(key string) (bool, error) {
	_, err := l.db.Get([]byte(key), nil)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Get ...
func (l *leveldb) Get(key string) ([]byte, error) {
	return l.db.Get([]byte(key), nil)
}

// Set ...
func (l *leveldb) Set(key string, value []byte) error {
	return l.db.Put([]byte(key), value, nil)
}

// Delete ...
func (l *leveldb) Delete(key string) error {
	return l.db.Delete([]byte(key), nil)
}

// Open ...
func (l *leveldb) Open() error {
	path, cleanTemp, err := files.TempDir()
	if err != nil {
		return err
	}

	db, err := level.OpenFile(path.Name(), nil)
	if err != nil {
		return err
	}
	l.db = db
	l.teardown = cleanTemp

	return nil
}

// Close ...
func (l *leveldb) Close() error {
	defer l.teardown()
	return l.db.Close()
}
