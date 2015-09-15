package store

import (
	"os"
	"sync"
)

// Store holds data, but relies on external services for data knowledge.
type Store struct {
	sync.Mutex

	file *os.File
}

// New data store object.
func New(filename string) (store *Store, err error) {

	// Verify the file can be opened.
	var file *os.File
	if file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666); nil != err {
		return
	}

	store = &Store{
		file: file,
	}
	return
}

// ReadAt returns the data at the requested location.
func (store *Store) ReadAt(b []byte, off uint64) (n uint64, err error) {
	store.Lock()
	defer store.Unlock()

	var cnt int
	cnt, err = store.file.ReadAt(b, int64(off))
	n = uint64(cnt)
	return
}

// Add an entry in the store.
func (store *Store) Add(b []byte) (off, length uint64, err error) {

	store.Lock()
	defer store.Unlock()

	var start int64
	if start, err = store.file.Seek(0, 2); nil != err {
		return
	}

	var n int
	if n, err = store.file.Write(b); nil != err {
		return
	}

	off = uint64(start)
	length = uint64(n)
	return
}

// Close the store.
func (store *Store) Close() (err error) {
	store.Lock()
	defer store.Unlock()

	return store.file.Close()
}
