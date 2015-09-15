package link

import (
	"encoding/binary"
	"errors"
	"os"
	"sync"
)

var (
	// ErrFailedToRead16Bytes is thrown when 16 bytes cannot be read from disk.
	ErrFailedToRead16Bytes = errors.New("Failed to read 16 bytes from disk")
	// ErrFailedToWrite16Bytes is thrown when 16 bytes cannot be written to disk.
	ErrFailedToWrite16Bytes = errors.New("Failed to write 16 bytes to disk")

	// empty for marking deletions, specifically the length.
	empty = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

// Link tracks data locations and lengths.
type Link struct {
	sync.Mutex

	file *os.File
}

// New Link object.
func New(filename string) (link *Link, err error) {

	// Verify the file can be opened.
	var file *os.File
	if file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666); nil != err {
		return
	}

	link = &Link{
		file: file,
	}
	return
}

// At returns the position and length of the requested data.
func (link *Link) At(id uint64) (start, length uint64, err error) {
	link.Lock()
	defer link.Unlock()

	if _, err = link.file.Seek(int64(id*16), 0); nil != err {
		return
	}

	entry := make([]byte, 16)
	if num, _ := link.file.Read(entry); 16 != num {
		err = ErrFailedToRead16Bytes
		return
	}

	start = binary.BigEndian.Uint64(entry[:8])
	length = binary.BigEndian.Uint64(entry[8:])
	return
}

// Add an entry in the link.
func (link *Link) Add(start, length uint64) (id uint64, err error) {
	entry := make([]byte, 16)
	binary.BigEndian.PutUint64(entry[:8], start)
	binary.BigEndian.PutUint64(entry[8:], length)

	link.Lock()
	defer link.Unlock()

	var pos int64
	if pos, err = link.file.Seek(0, 2); nil != err {
		return
	}

	var n int
	if n, err = link.file.Write(entry); nil != err {
		return
	} else if 16 != n {
		err = ErrFailedToWrite16Bytes
		return
	}

	id = uint64(pos / 16)
	return
}

// Delete an entry in the link.
func (link *Link) Delete(id uint64) (freed uint64, err error) {
	link.Lock()
	defer link.Unlock()

	if _, err = link.file.Seek(int64(id*16+8), 0); nil != err {
		return
	}

	entry := make([]byte, 8)
	if num, _ := link.file.Read(entry); 8 != num {
		err = ErrFailedToRead16Bytes // Actually 8 bytes, here.
		return
	}
	freed = binary.BigEndian.Uint64(entry)

	if _, err = link.file.Seek(int64(id*16), 0); nil != err {
		return
	}

	var n int
	if n, err = link.file.Write(empty); nil != err {
		return
	} else if 16 != n {
		err = ErrFailedToWrite16Bytes
		return
	}
	return
}

// Close the link.
func (link *Link) Close() (err error) {
	link.Lock()
	defer link.Unlock()

	return link.file.Close()
}
