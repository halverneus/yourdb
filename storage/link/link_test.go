package link

import (
	"os"
	"testing"
)

const (
	testfile = "mylink.file"
)

var (
	data []*record
)

type record struct {
	start, length, id uint64
}

func TestMain(m *testing.M) {
	buildUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func buildUp() {
	start := uint64(1000)
	length := uint64(1000000)
	for i := 0; i < 100; i++ {
		data = append(data, &record{start, length, uint64(0)})
		start += 100
		length -= 100
	}

	os.Remove(testfile)
}

func tearDown() {
	os.Remove(testfile)
}

func TestFirstLoadAndClose(t *testing.T) {
	link, err := New(testfile)
	if nil != err {
		t.Error(err)
	}
	if err = link.Close(); nil != err {
		t.Error(err)
	}
}

func TestSecondLoadAndClose(t *testing.T) {
	link, err := New(testfile)
	if nil != err {
		t.Error(err)
	}
	if err = link.Close(); nil != err {
		t.Error(err)
	}
}

func TestReadFromEmpty(t *testing.T) {
	link, err := New(testfile)
	if nil != err {
		t.Error(err)
	}

	if _, _, err = link.At(0); nil == err {
		t.Error("Should have failed to read id 0.")
	}
	if _, _, err = link.At(10); nil == err {
		t.Error("Should have failed to read id 10.")
	}

	if err = link.Close(); nil != err {
		t.Error(err)
	}
}

func TestWriteAndReadOne(t *testing.T) {
	link, err := New(testfile)
	if nil != err {
		t.Error(err)
	}

	id, err := link.Add(10, 20)
	if nil != err {
		t.Error(err)
	} else if 0 != id {
		t.Error("Expected to write first entry.")
	}

	start, length, err := link.At(id)
	if nil != err {
		t.Error(err)
	}
	if 10 != start {
		t.Error("Expected start of 10, got:", start)
	}
	if 20 != length {
		t.Error("Expected length of 20, got:", length)
	}

	if err = link.Close(); nil != err {
		t.Error(err)
	}
}

func TestWriteMore(t *testing.T) {

	link, err := New(testfile)
	if nil != err {
		t.Error(err)
	}

	for _, item := range data {
		if item.id, err = link.Add(item.start, item.length); nil != err {
			t.Error(err)
		}
	}

	for _, item := range data {
		st, le, err := link.At(item.id)
		if nil != err {
			t.Error(err)
		}
		if st != item.start {
			t.Error("Expected start of", item.start, ", got", st)
		}
		if le != item.length {
			t.Error("Expected length of", item.length, ", got", le)
		}
	}

	if err = link.Close(); nil != err {
		t.Error(err)
	}
}

func TestReadMore(t *testing.T) {
	link, err := New(testfile)
	if nil != err {
		t.Error(err)
	}

	for _, item := range data {
		st, le, err := link.At(item.id)
		if nil != err {
			t.Error(err)
		}
		if st != item.start {
			t.Error("Expected start of", item.start, ", got", st)
		}
		if le != item.length {
			t.Error("Expected length of", item.length, ", got", le)
		}
	}

	if err = link.Close(); nil != err {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	link, err := New(testfile)
	if nil != err {
		t.Error(err)
	}

	periodic := 0
	for _, item := range data {
		if 5 == periodic {
			periodic = 0
			var length uint64
			if length, err = link.Delete(item.id); nil != err {
				t.Error(err)
			}
			if length != item.length {
				t.Error("Expected to free", item.length, "bytes and freed", length)
			}
			item.start = 0
			item.length = 0
		} else {
			periodic++
		}
	}

	for _, item := range data {
		st, le, err := link.At(item.id)
		if nil != err {
			t.Error(err)
		}
		if st != item.start {
			t.Error("Expected start of", item.start, ", got", st)
		}
		if le != item.length {
			t.Error("Expected length of", item.length, ", got", le)
		}
	}

	if err = link.Close(); nil != err {
		t.Error(err)
	}
}
