package gost_test

import (
	"github.com/egorgasay/gost"
	"testing"
)

func TestResult(t *testing.T) {
	f := func() gost.Result[int] {
		return gost.Ok(1)
	}

	if res := f().Unwrap(); res != 1 {
		t.Fatal("unexpected value:", res)
	}

	f = func() gost.Result[int] {
		return gost.Err[int](gost.NewError(1, 2, "test"))
	}

	if res := f().Error(); res.BaseCode() != 1 || res.ExtendedCode() != 2 || res.Message() != "test" {
		t.Fatal("unexpected error:", res)
	}

	_ = func() gost.ResultN {
		return gost.ErrN(gost.NewError(1, 2, "test"))
	}
}
