package storage

import (
	"bytes"
	"crypto/rand"
	"io"
	"os"
	"testing"
)

const (
	testfile = "mystorage"
)

var (
	data []*record
)

type record struct {
	id       uint64
	contents []byte
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

func gen1K() (data []byte) {
	for i := 0; i < 64; i++ {
		rd := gen()
		data = append(data, rd...)
	}
	return
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
		data = append(data, &record{uint64(0), contents})
	}

	os.Remove(testfile + linkExt)
	os.Remove(testfile + metaExt)
	os.Remove(testfile + storeExt)
}

func tearDown() {
	os.Remove(testfile + linkExt)
	os.Remove(testfile + metaExt)
	os.Remove(testfile + storeExt)
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
		item.id, err = store.Add(item.contents)
		if nil != err {
			t.Error(err)
		}
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
		fresh, err := store.Get(item.id)
		if nil != err {
			t.Error(err)
		}

		if 0 != bytes.Compare(fresh, item.contents) {
			t.Error("Data does not match.")
		}
	}

	if err = store.Close(); nil != err {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	store, err := New(testfile)
	if nil != err {
		t.Error(err)
	}

	for _, item := range data {
		if err = store.Delete(item.id); nil != err {
			t.Error(err)
		}
	}

	for _, item := range data {
		fresh, err := store.Get(item.id)
		if nil != err {
			t.Error(err)
		}
		if 0 < len(fresh) {
			t.Error("Expected null returned")
		}
	}

	if err = store.Close(); nil != err {
		t.Error(err)
	}
}

func TestWipe(t *testing.T) {
	store, err := New(testfile)
	if nil != err {
		t.Error(err)
	}

	if err = store.Close(); nil != err {
		t.Error(err)
	}

	if err = store.Wipe(); nil != err {
		t.Error(err)
	}
}

func BenchmarkWrite(b *testing.B) {
	mydata := gen1K()
	store, _ := New(testfile)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Add(mydata)
	}

	store.Close()
}

func BenchmarkRead(b *testing.B) {
	mydata := gen1K()
	store, _ := New(testfile)

	id, _ := store.Add(mydata)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Get(id)
	}

	store.Close()
}

func BenchmarkDelete(b *testing.B) {
	mydata := gen1K()
	store, _ := New(testfile)

	for i := 0; i < b.N; i++ {
		store.Add(mydata)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store.Delete(uint64(i))
	}

	store.Close()
}
