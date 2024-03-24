package gost

import (
	"errors"
	"fmt"
	"testing"
)

func TestError_WrapWithMessage(t *testing.T) {
	e := &Error{baseCode: 1, extendedCode: 2, message: "3"}

	e = e.Wrap("test")

	if e.Message() != "test: 3" {
		t.Fatal("unexpected error:", e)
	}

	e = nil
	e = e.Wrap("test")

	if e != nil {
		t.Fatal("unexpected error:", e)
	}
}

func TestErrorIs(t *testing.T) {

	oneErr := NewErrX(1, "3")

	secErr := oneErr.Extend(2, "4")

	//if !secErr.Is(oneErr) {
	//	t.Fatal("unexpected error:", secErr)
	//}

	if !errors.Is(secErr, oneErr) {
		t.Fatal("unexpected error:", secErr)
	}

	//if secErr.Is(thirdErr) {
	//	t.Fatal("unexpected error:", secErr)
	//}
	//
	//fourthErr := NewError(1, 2, "3").Wrap("hi").Join(nil).Join(nil)
	//
	//if fourthErr.Is(nil) {
	//	t.Fatal("unexpected error:", fourthErr)
	//}
}

func TestErrXIs(t *testing.T) {
	var cases = []struct {
		name  string
		build func() (err error, target error)
	}{
		{"same", func() (err error, target error) {
			err1 := NewErrX(1, "3")

			return err1, err1
		}},
		{"extWith_ExtendMsg", func() (err error, target error) {
			err1 := NewErrX(1, "3")

			return err1, err1.ExtendMsg("12")
		}},
		{"extWith_Extend", func() (err error, target error) {
			err1 := NewErrX(1, "3")

			return err1, err1.Extend(2, "12")
		}},
		{"extWith_Extend2", func() (err error, target error) {
			err1 := NewErrX(1, "3")

			return err1, err1.
				Extend(2, "12").
				Extend(3, "13").
				Extend(4, "14").
				Extend(5, "15")
		}},
		{"extWith_Extend_fmt_Errorf", func() (err error, target error) {
			err1 := NewErrX(1, "3")

			return err1, fmt.Errorf("%w: %w", err1.Extend(2, "12"), err1)
		}},
		{"extWith_Extend_errors_Join", func() (err error, target error) {
			err1 := NewErrX(1, "3")
			err2 := NewErrX(2, "4")

			return err1, errors.Join(err1, err2)
		}},
		{"extWith_Extend_errors_Join#2", func() (err error, target error) {
			err1 := NewErrX(1, "3")
			err2 := NewErrX(2, "4").Extend(3, "5")

			return err1, errors.Join(err1, err2)
		}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err, target := c.build()
			if !errors.Is(err, target) {
				t.Fatal("unexpected error:", err)
			}
		})
	}
}
