package gost_test

import (
	"github.com/egorgasay/gost"
	"testing"
)

type Err interface {
	err()
}

type (
	ErrNotFound  struct{}
	ErrForbidden struct{}
)

func (ErrNotFound) err()  {}
func (ErrForbidden) err() {}

func TestEnumComplicated(t *testing.T) {
	e := func() (enum gost.Enum[Err]) {
		return enum.New(ErrForbidden{})
	}()

	switch v := e.Enum().(type) {
	case ErrNotFound:
		t.Fatalf("Error, wrong case: %v", v)
	case ErrForbidden:
		t.Log("Ok")
	default:
		t.Fatal("Error, wrong type:", v)
	}
}
