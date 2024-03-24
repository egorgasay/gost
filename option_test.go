package gost_test

import (
	"testing"

	"github.com/egorgasay/gost"
)

func TestOption(t *testing.T) {
	f := func() gost.Option[int] {
		return gost.Some(1)
	}

	if res := f().Unwrap(); res != 1 {
		t.Fatal("unexpected value:", res)
	}

	f = func() gost.Option[int] {
		return gost.None[int]()
	}

	if res := f().UnwrapOrDefault(); res != 0 {
		t.Fatal("unexpected value:", res)
	}
}
