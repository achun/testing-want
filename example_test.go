package want_test

import (
	"errors"
	"github.com/achun/testing-want"
	"testing"
)

func TestWants(t *testing.T) {
	ExampleT(t)
	ExampleCaller(t)
	ExampleTrue(t)
	ExampleEqual(t)
	ExampleRecover(t)
	ExamplePanic(t)
	ExampleNil(t)
	ExampleError(t)
}

func ExampleCaller(t *testing.T) {
	want.Equal(t, want.Caller(1), "github.com/achun/testing-want_test.ExampleCaller:21")
	want.Equal(t, want.Caller(2), "github.com/achun/testing-want_test.TestWants:11")
}

func ExampleT(t *testing.T) {
	wt := want.T(t)

	wt.Equal(want.Caller(1), "github.com/achun/testing-want_test.ExampleT:28")
	wt.Equal(want.Caller(2), "github.com/achun/testing-want_test.TestWants:10")
}

func ExampleTrue(t *testing.T) {
	want.True(t, 1 == 1)
}

func ExampleEqual(t *testing.T) {
	wants := "something"
	got := want.String(wants)
	want.Equal(t, got, "\n"+wants, "...")
}

func ExampleRecover(t *testing.T) {
	want.Recover(t, "<nil>", func() {
		// your code
		// panic("Are you sure?")
	})
}

func ExamplePanic(t *testing.T) {
	want.Panic(t, func() {
		panic("Are you sure?")
	})
}

func ExampleNil(t *testing.T) {
	want.Nil(t, nil, "...")
}

func ExampleError(t *testing.T) {
	want.Error(t, errors.New("newError"), "...")
}
