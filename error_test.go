package gost

import (
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

	oneErr := NewError(1, 2, "3")

	secErr := NewError(0, 0, "Sec").Join(oneErr)

	thirdErr := NewError(1, 2, "3")

	if !secErr.Is(oneErr) {
		t.Fatal("unexpected error:", secErr)
	}

	if secErr.Is(thirdErr) {
		t.Fatal("unexpected error:", secErr)
	}

	fourthErr := NewError(1, 2, "3").Wrap("hi").Join(nil).Join(nil)

	if fourthErr.Is(nil) {
		t.Fatal("unexpected error:", fourthErr)
	}
}
