package gost

import (
	"testing"
)

func TestError_WrapWithMessage(t *testing.T) {
	e := &Error{
		baseCode:     1,
		extendedCode: 2,
		message:      "3",
	}

	e = e.WrapNotNilMsg("test")

	if e.Message() != "test: 3" {
		t.Fatal("unexpected error:", e)
	}

	e = nil

	e = e.WrapNotNilMsg("test")

	if e != nil {
		t.Fatal("unexpected error:", e)
	}
}
