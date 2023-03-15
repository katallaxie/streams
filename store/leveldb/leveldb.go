package leveldb

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ionos-cloud/streams/store"

	"github.com/katallaxie/pkg/utils/files"
	level "github.com/syndtr/goleveldb/leveldb"
)

type leveldb struct {
	db       *level.DB
	teardown files.TearDownFunc

	store.Unimplemented
}

// New creates a new leveldb storage.
func New() store.Storage {
	l := new(leveldb)

	return l
}

// Has is called to check if a key exists.
func (l *leveldb) Has(key string) (bool, error) {
	_, err := l.db.Get([]byte(key), nil)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Get is called to get a value.
func (l *leveldb) Get(key string) ([]byte, error) {
	return l.db.Get([]byte(key), nil)
}

// Set is called to set a value.
func (l *leveldb) Set(key string, value []byte) error {
	return l.db.Put([]byte(key), value, nil)
}

// Delete is called to delete a value.
func (l *leveldb) Delete(key string) error {
	return l.db.Delete([]byte(key), nil)
}

// Open is called when the storage is opened.
func (l *leveldb) Open() error {
	name := time.Now().Format("20060102150405")

	db, err := level.OpenFile(filepath.Join(os.TempDir(), name), nil)
	if err != nil {
		return err
	}

	l.db = db
	l.teardown = func() {
		_ = os.Remove(filepath.Join(os.TempDir(), name))
	}

	return nil
}

// Close is called when the storage is closed.
func (l *leveldb) Close() error {
	defer l.teardown()

	return l.db.Close()
}
