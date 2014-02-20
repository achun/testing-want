package want

import (
	"fmt"
	"runtime"
	"testing"
)

// returns string by fmt.Sprint. if not empty, prefix "\n"
func String(show ...interface{}) string {
	if len(show) != 0 {
		return "\n" + fmt.Sprint(show...)
	}
	return ""
}

// returns filepath and line by runtime
func Caller(skip int) string {
	pc, _, line, _ := runtime.Caller(skip)
	return runtime.FuncForPC(pc).Name() + ":" + fmt.Sprint(line) + "\n"
}

// want ok equal true
func True(t *testing.T, ok bool, show ...interface{}) {
	if !ok {
		t.Fatal(Caller(2), String(show...))
	}
}

// want got equal wants
func Equal(t *testing.T, got, wants interface{}, show ...interface{}) {
	if wants != got {
		t.Fatal(Caller(2), "want:", wants, ", but got:", got, String(show...))
	}
}

// want recover msg string
func Recover(t *testing.T, msg string, fn func()) {
	defer func() {
		str := fmt.Sprint(recover())
		if msg != str {
			t.Fatal(Caller(3), "want recover:", msg, ", but got:", str)
		}
	}()
	fn()
}

// want recover panic
func Panic(t *testing.T, fn func()) {
	defer func() {
		if nil == recover() {
			t.Fatal(Caller(3), "want panic, but got nil")
		}
	}()
	fn()
}

// wants equal nil
func Nil(t *testing.T, wants interface{}, show ...interface{}) {
	if wants != nil {
		t.Fatal(Caller(2), "want nil, but got:", wants, String(show...))
	}
}

// wants not nil
func Error(t *testing.T, wants error, show ...interface{}) {
	if wants == nil {
		t.Fatal(Caller(2), "want an error but got nil.", String(show...))
	}
}
