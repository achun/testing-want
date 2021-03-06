package want

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func Println(i ...interface{}) {
	fmt.Println(i...)
}

func Printf(format string, i ...interface{}) {
	fmt.Printf(format, i...)
}

// returns string by fmt.Sprint. if not empty, prefix "\n"
func String(show ...interface{}) string {
	if len(show) == 0 {
		return ""
	}
	ret := "\n"
	for _, i := range show {
		fn, ok := i.(func() string)
		if ok {
			ret += fn()
		} else {
			ret += fmt.Sprint(i)
		}
	}
	return ret
}

var LocalFileLine bool

// returns filepath and line by runtime
func Caller(skip int) string {
	pc, file, line, _ := runtime.Caller(skip)
	if LocalFileLine {
		return fmt.Sprintf("\n%s:\n  %v\n", file, line)
	}
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

func asPtr(k reflect.Kind) bool {
	switch k {
	case reflect.Chan,
		reflect.Func,
		reflect.Map,
		reflect.Ptr,
		reflect.Interface,
		reflect.Slice:

		return true
	}
	return false
}

func T(t testing.TB) Want {
	return Want{t, 2}
}

type Want struct {
	T    testing.TB
	Skip int
}

func (w Want) Fatalf(format string, show ...interface{}) {
	show = append([]interface{}{
		Caller(w.Skip) + "\n",
	}, show...)
	w.T.Fatalf(format, show...)
}

func (w Want) Fatal(show ...interface{}) {
	show = append([]interface{}{
		Caller(w.Skip) + "\n",
	}, show...)
	w.T.Fatal(show...)
}

func (w Want) True(ok bool, show ...interface{}) Want {
	if !ok {
		w.T.Fatal(Caller(w.Skip), "\nwant: true\n got: false", String(show...))
	}
	return w
}

func (w Want) False(ok bool, show ...interface{}) Want {
	if ok {
		w.T.Fatal(Caller(w.Skip), "\nwant: false\n got: true", String(show...))
	}
	return w
}

func (w Want) Equal(got, wants interface{}, show ...interface{}) Want {
	ser, ok := got.(fmt.Stringer)
	if ok {
		got = ser.String()
	} else if v := reflect.ValueOf(got); v.Kind() >= reflect.Array {
		got = fmt.Sprint(got)
	}

	ser, ok = wants.(fmt.Stringer)
	if ok {
		wants = ser.String()
	} else if v := reflect.ValueOf(wants); v.Kind() >= reflect.Array {
		wants = fmt.Sprint(wants)
	}

	if wants != got {
		w.T.Fatal(Caller(w.Skip), "\nwant:", wants, "\n got:", got, String(show...))
	}
	return w
}

func (w Want) Recover(msg string, fn func()) Want {
	defer func() {
		str := fmt.Sprint(recover())
		if msg != str {
			w.T.Fatal(Caller(w.Skip+1), "\nwant: recover ", msg, "\n got:", str)
		}
	}()
	fn()
	return w
}

func (w Want) Panic(fn func()) Want {
	defer func() {
		if nil == recover() {
			w.T.Fatal(Caller(w.Skip+1), "\nwant: panic\n got: nil")
		}
	}()
	fn()
	return w
}

func (w Want) Nil(i interface{}, show ...interface{}) Want {
	v := reflect.ValueOf(i)
	if asPtr(v.Kind()) && !v.IsNil() || v.IsValid() {
		w.T.Fatal(Caller(w.Skip), "\nwant: nil\n got:", i, String(show...))
	}
	return w
}

func (w Want) NotNil(i interface{}, show ...interface{}) Want {
	v := reflect.ValueOf(i)
	if !v.IsValid() || asPtr(v.Kind()) && v.IsNil() {
		w.T.Fatal(Caller(w.Skip), "\nwant: not nil\n got:", i, String(show...))
	}
	return w
}

func (w Want) Error(err error, show ...interface{}) Want {
	if err == nil {
		w.T.Fatal(Caller(w.Skip), "\nwant: an error\n got: nil", String(show...))
	}
	return w
}

// want ok equal true
func True(t *testing.T, ok bool, show ...interface{}) {
	(Want{t, 3}).True(ok, show...)
}

// want ok equal false
func False(t *testing.T, ok bool, show ...interface{}) {
	(Want{t, 3}).False(ok, show...)
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

// want v as nil, if not nil Fatal
func Nil(t *testing.T, v interface{}, show ...interface{}) {
	(Want{t, 3}).Nil(v, show...)
}

// want v not nil, if not nil Fatal
func NotNil(t *testing.T, v interface{}, show ...interface{}) {
	(Want{t, 3}).NotNil(v, show...)
}

// wants an error, if nil Fatal
func Error(t *testing.T, err error, show ...interface{}) {
	(Want{t, 3}).Error(err, show...)
}
