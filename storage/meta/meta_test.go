package meta

import (
	"os"
	"testing"
)

const (
	testfile = "mymeta.file"
)

func TestMain(m *testing.M) {
	buildUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func buildUp() {

}

func tearDown() {
	os.Remove(testfile)
}

func TestFirstLoadAndClose(t *testing.T) {
	if meta, err := New(testfile); nil != err {
		t.Error(err)
	} else if nil == meta {
		t.Error("No error, no object.")
	} else {
		if err = meta.Close(); nil != err {
			t.Error(err)
		}
	}
}

func TestSecondLoadAndClose(t *testing.T) {
	if meta, err := New(testfile); nil != err {
		t.Error(err)
	} else if nil == meta {
		t.Error("No error, no object.")
	} else {
		if err = meta.Close(); nil != err {
			t.Error(err)
		}
	}
}

func TestNext(t *testing.T) {
	meta, err := New(testfile)
	if nil != err {
		t.Error(err)
	}

	if pos, err := meta.Next(99); nil != err {
		t.Error(err)
	} else if 0 != pos {
		t.Error("Expected to receive beginning of file.")
	}
	if pos, err := meta.Next(4); nil != err {
		t.Error(err)
	} else if 99 != pos {
		t.Error("Expected position of 99")
	}

	if err = meta.Close(); nil != err {
		t.Error(err)
	}
}

func TestNextAgain(t *testing.T) {
	meta, err := New(testfile)
	if nil != err {
		t.Error(err)
	}

	if pos, err := meta.Next(30); nil != err {
		t.Error(err)
	} else if 103 != pos {
		t.Error("Expected to receive beginning of file.")
	}
	if pos, err := meta.Next(4); nil != err {
		t.Error(err)
	} else if 133 != pos {
		t.Error("Expected position of 99")
	}

	if err = meta.Close(); nil != err {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	meta, err := New(testfile)
	if nil != err {
		t.Error(err)
	}

	if total, err := meta.Deleted(0); nil != err {
		t.Error(err)
	} else if 0 != total {
		t.Error("Expected to have 0 deleted bytes, got:", total)
	}

	if total, err := meta.Deleted(30); nil != err {
		t.Error(err)
	} else if 30 != total {
		t.Error("Expected to have 30 deleted bytes, got:", total)
	}

	if err = meta.Close(); nil != err {
		t.Error(err)
	}
}

func TestDeleteAgain(t *testing.T) {
	meta, err := New(testfile)
	if nil != err {
		t.Error(err)
	}

	if total, err := meta.Deleted(0); nil != err {
		t.Error(err)
	} else if 30 != total {
		t.Error("Expected to have 0 deleted bytes, got:", total)
	}

	if total, err := meta.Deleted(40); nil != err {
		t.Error(err)
	} else if 70 != total {
		t.Error("Expected to have 30 deleted bytes, got:", total)
	}

	if total, err := meta.Deleted(90); nil != err {
		t.Error(err)
	} else if 160 != total {
		t.Error("Expected to have 30 deleted bytes, got:", total)
	}

	if err = meta.Close(); nil != err {
		t.Error(err)
	}
}
