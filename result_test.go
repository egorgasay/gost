package gost_test

import (
	"testing"

	"github.com/egorgasay/gost"
)

func TestResult(t *testing.T) {
	f := func() gost.Result[int] {
		return gost.Ok(1)
	}

	if res := f().Unwrap(); res != 1 {
		t.Fatal("unexpected value:", res)
	}

	f = func() gost.Result[int] {
		return gost.Err[int](gost.NewErrX(1, "test").Extend(2))
	}

	if res := f().Error(); res.BaseCode() != 1 || res.ExtCode() != 2 || res.MessagesSpace() != "test" {
		t.Fatal("unexpected error:", res)
	}

	_ = func() gost.ResultN {
		return gost.ErrN(gost.NewErrX(1, "test").Extend(2))
	}().Error()
}
