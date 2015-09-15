package storage

import (
	"errors"
	"github.com/halverneus/yourdb/storage/link"
	"github.com/halverneus/yourdb/storage/meta"
	"github.com/halverneus/yourdb/storage/store"
	"os"
)

var (
	// ErrPartialData is the result of an unexpected amount of data being read.
	ErrPartialData = errors.New("Partial data recovered")
)

const (
	linkExt  = ".link"
	metaExt  = ".meta"
	storeExt = ".store"
)

// Storage for objects.
type Storage struct {
	base  string
	link  *link.Link
	meta  *meta.Meta
	store *store.Store
}

// New storage constructor.
func New(baseFilename string) (storage *Storage, err error) {
	storage = &Storage{
		base: baseFilename,
	}
	if storage.link, err = link.New(storage.base + linkExt); nil != err {
		return
	}
	if storage.meta, err = meta.New(storage.base + metaExt); nil != err {
		return
	}
	if storage.store, err = store.New(storage.base + storeExt); nil != err {
		return
	}
	return
}

// Add data to storage.
func (storage *Storage) Add(data []byte) (id uint64, err error) {
	start := uint64(0)
	length := uint64(0)
	if start, length, err = storage.store.Add(data); nil != err {
		return
	}
	if id, err = storage.link.Add(start, length); nil != err {
		return
	}
	if _, err = storage.meta.Next(length); nil != err {
		return
	}

	return
}

// Get data stored at key 'id'.
func (storage *Storage) Get(id uint64) (data []byte, err error) {
	start := uint64(0)
	length := uint64(0)
	if start, length, err = storage.link.At(id); nil != err {
		return
	}
	if 0 == length {
		// No error. Value is just nil.
		return
	}

	data = make([]byte, length)

	read := uint64(0)
	if read, err = storage.store.ReadAt(data, start); nil != err {
		return
	}
	if read != length {
		err = ErrPartialData
		return
	}
	return
}

// Delete data stored at key 'id'.
func (storage *Storage) Delete(id uint64) (err error) {
	length := uint64(0)
	if length, err = storage.link.Delete(id); nil != err {
		return
	}
	if 0 < length {
		if _, err = storage.meta.Deleted(length); nil != err {
			return
		}
	}
	return
}

// Close all file connections.
func (storage *Storage) Close() (err error) {
	if err = storage.store.Close(); nil != err {
		storage.link.Close()
		storage.meta.Close()
		return
	}
	if err = storage.link.Close(); nil != err {
		storage.meta.Close()
		return
	}
	err = storage.meta.Close()
	return
}

// Wipe all storage files from disk.
func (storage *Storage) Wipe() (err error) {
	if err = os.Remove(storage.base + storeExt); nil != err {
		os.Remove(storage.base + metaExt)
		os.Remove(storage.base + linkExt)
		return
	}
	if err = os.Remove(storage.base + metaExt); nil != err {
		os.Remove(storage.base + linkExt)
		return
	}
	err = os.Remove(storage.base + linkExt)
	return
}
