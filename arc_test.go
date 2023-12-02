package gost_test

import (
	"github.com/egorgasay/gost"
	"testing"
)

func TestArc(t *testing.T) {
	type S struct {
		x int
		y string
	}

	s := S{
		x: 1,
		y: "test",
	}

	ref := gost.NewArc(&s)

	if ref.Read().x != 1 {
		t.Fatal("unexpected value:", ref.Read().x)
	}

	protected := ref.Read()
	protected.x = 2

	if ref.Read().x == 2 {
		t.Fatal("unexpected value: 2")
	}

	if ref.Read().y != "test" {
		t.Fatal("unexpected value:", ref.Read().y)
	}

	protected2 := ref.Read()
	protected2.y = "test2"

	if ref.Read().y == "test2" {
		t.Fatal("unexpected value: test2")
	}
}
