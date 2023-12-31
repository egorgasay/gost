package gost_test

import (
	"testing"
	"time"

	"github.com/egorgasay/gost"
)

func TestMutexOk(t *testing.T) {
	var protected gost.Mutex[int] = gost.NewMutex(1)

	for i := 0; i < 10; i++ {
		go func() {
			*protected.BorrowMut()++
			protected.Release()
		}()
	}

	time.Sleep(1 * time.Second)

	if protected.Borrow().Read() != 11 {
		t.Fatal()
	}
}

func TestMutexErr(t *testing.T) {
	var protected gost.Mutex[int] = gost.NewMutex(1)

	*protected.BorrowMut()++

	err := protected.TryBorrow().Error()
	if err != gost.ErrMutexIsLocked {
		t.Fatal("unexpected error:", err)
	}
}

func TestSecureMutex(t *testing.T) {
	var prot = gost.NewSecureMutex[int]().Set(2).Lock()

	if prot.Unlock().Value().Read() != 2 {
		t.Fatal()
	}
}
