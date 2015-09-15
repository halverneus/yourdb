package store

import (
	"bytes"
	"crypto/rand"
	"io"
	"os"
	"testing"
)

const (
	testfile = "mystore.file"
)

var (
	data []*record
)

type record struct {
	start, length uint64
	contents      []byte
}

func gen() []byte {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n == len(uuid) && err == nil {
		// variant bits; see section 4.1.1
		uuid[8] = uuid[8]&^0xc0 | 0x80
		// version 4 (pseudo-random); see section 4.1.3
		uuid[6] = uuid[6]&^0xf0 | 0x40
	}
	return uuid
}

func TestMain(m *testing.M) {
	buildUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func buildUp() {
	for i := 0; i < 100; i++ {
		var contents []byte
		for j := 1; j < i; j++ {
			rd := gen()
			contents = append(contents, rd...)
		}
		data = append(data, &record{uint64(0), uint64(len(contents)), contents})
	}

	os.Remove(testfile)
}

func tearDown() {
	os.Remove(testfile)
}

func TestFirstLoadAndClose(t *testing.T) {
	store, err := New(testfile)
	if nil != err {
		t.Error(err)
	}
	if err = store.Close(); nil != err {
		t.Error(err)
	}
}

func TestSecondLoadAndClose(t *testing.T) {
	store, err := New(testfile)
	if nil != err {
		t.Error(err)
	}
	if err = store.Close(); nil != err {
		t.Error(err)
	}
}

func TestWrites(t *testing.T) {
	store, err := New(testfile)
	if nil != err {
		t.Error(err)
	}

	for _, item := range data {
		start, length, err := store.Add(item.contents)
		if nil != err {
			t.Error(err)
		}
		if length != item.length {
			t.Error("Wrong number of bytes written.")
		}
		item.start = start
	}

	if err = store.Close(); nil != err {
		t.Error(err)
	}
}

func TestReads(t *testing.T) {
	store, err := New(testfile)
	if nil != err {
		t.Error(err)
	}

	for _, item := range data {
		fresh := make([]byte, item.length)
		length, err := store.ReadAt(fresh, item.start)
		if nil != err {
			t.Error(err)
		}
		if length != item.length {
			t.Error("Expected length of", item.length, ", got", length)
		}

		if 0 != bytes.Compare(fresh, item.contents) {
			t.Error("Data does not match.")
		}
	}

	if err = store.Close(); nil != err {
		t.Error(err)
	}
}
