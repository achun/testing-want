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
	return runtime.FuncForPC(pc).Name() + ":" + fmt.Sprint(line)
}

// return last argument.(error)
func LastError(rets ...interface{}) error {
	l := len(rets)
	if l == 0 {
		return nil
	}
	err, _ := rets[l-1].(error)
	return err
}

func T(t *testing.T) Want {
	return Want{t, 2}
}

type Want struct {
	T    *testing.T
	Skip int
}

func (w Want) True(ok bool, show ...interface{}) Want {
	if !ok {
		w.T.Fatal(Caller(w.Skip), "\n", String(show...))
	}
	return w
}

func (w Want) False(ok bool, show ...interface{}) Want {
	if ok {
		w.T.Fatal(Caller(w.Skip), "\n", String(show...))
	}
	return w
}

func (w Want) Equal(got, wants interface{}, show ...interface{}) Want {
	if wants != got {
		w.T.Fatal(Caller(w.Skip), "\nwant:", wants, ", but got:", got, String(show...))
	}
	return w
}

func (w Want) Recover(msg string, fn func()) Want {
	defer func() {
		str := fmt.Sprint(recover())
		if msg != str {
			w.T.Fatal(Caller(w.Skip+1), "\nwant recover:", msg, ", but got:", str)
		}
	}()
	fn()
	return w
}

func (w Want) Panic(fn func()) Want {
	defer func() {
		if nil == recover() {
			w.T.Fatal(Caller(w.Skip+1), "\nwant panic, but got nil")
		}
	}()
	fn()
	return w
}

func (w Want) Nil(err interface{}, show ...interface{}) Want {
	if err != nil {
		w.T.Fatal(Caller(w.Skip), "\nwant nil, but got:", err, String(show...))
	}
	return w
}

func (w Want) Error(err error, show ...interface{}) Want {
	if err == nil {
		w.T.Fatal(Caller(w.Skip), "\nwant an error but got nil.", String(show...))
	}
	return w
}

// want ok equal true
func True(t *testing.T, ok bool, show ...interface{}) {
	(Want{t, 3}).True(ok, show...)
}

// want got equal wants
func Equal(t *testing.T, got, wants interface{}, show ...interface{}) {
	(Want{t, 3}).Equal(got, wants, show...)
}

// want recover msg string
func Recover(t *testing.T, msg string, fn func()) {
	(Want{t, 3}).Recover(msg, fn)
}

// want recover panic, if nil Fatal
func Panic(t *testing.T, fn func()) {
	(Want{t, 3}).Panic(fn)
}

// want err equal nil, if not nil Fatal
func Nil(t *testing.T, err interface{}, show ...interface{}) {
	(Want{t, 3}).Nil(err, show...)
}

// wants an error, if nil Fatal
func Error(t *testing.T, err error, show ...interface{}) {
	(Want{t, 3}).Error(err, show...)
}
