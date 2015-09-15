package meta

import (
	"encoding/gob"
	"os"
	"sync"
)

// Meta data used to track file state.
type Meta struct {
	sync.Mutex

	file *os.File
	data *data
}

type data struct {
	Deleted uint64
	End     uint64
}

// New Meta object.
func New(filename string) (meta *Meta, err error) {

	// Verify the file can be opened.
	var file *os.File
	if file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666); nil != err {
		return
	}

	// Create the meta object.
	meta = &Meta{
		file: file,
		data: &data{},
	}

	// Attempt to load the meta data...
	if err = meta.load(); nil != err {

		// If it fails to load, attempt to initialize as save the meta data...
		meta.data.Deleted = 0
		meta.data.End = 0
		if err = meta.save(); nil != err {

			// Initialization failed.
			meta = nil
			return
		}
	}

	return
}

// Next pointer in file for writing [n] number of bytes.
func (meta *Meta) Next(n uint64) (pos uint64, err error) {
	meta.Lock()
	defer meta.Unlock()

	pos = meta.data.End
	meta.data.End += n
	err = meta.save()
	return
}

// Deleted [n] number of bytes. Returns current total.
func (meta *Meta) Deleted(n uint64) (total uint64, err error) {
	meta.Lock()
	defer meta.Unlock()

	meta.data.Deleted += n
	total = meta.data.Deleted

	if n > 0 {
		err = meta.save()
	}
	return
}

// Close the meta data.
func (meta *Meta) Close() (err error) {
	meta.Lock()
	defer meta.Unlock()

	return meta.file.Close()
}

// Load the meta data into memory.
func (meta *Meta) load() (err error) {

	// Go to the beginning of the meta data file.
	if _, err = meta.file.Seek(0, 0); nil != err {
		return
	}

	// Decode the meta data into memory.
	decoder := gob.NewDecoder(meta.file)
	err = decoder.Decode(meta.data)
	return
}

// Save the meta data to disk.
func (meta *Meta) save() (err error) {

	// Go to the beginning of the meta data file to overwrite.
	if _, err = meta.file.Seek(0, 0); nil != err {
		return
	}

	// Encode the meta data to disk.
	encoder := gob.NewEncoder(meta.file)
	err = encoder.Encode(meta.data)
	return
}
