package gost

import (
	"fmt"
	"sync"
)

type scopedLock struct {
	mu sync.Mutex
}

type debugScopedLock struct {
	mu  sync.Mutex
	log func(isAboutToLock bool)
}

func (d *debugScopedLock) Lock() func() {
	d.log(true)
	d.mu.Lock()
	return func() {
		d.log(false)
		d.mu.Unlock()
	}
}

type Scoped interface {
	Lock() func()
}

func NewScopedLock() Scoped {
	return &scopedLock{}
}

func NewDebugScopedLock(log ...func(isAboutToLock bool)) Scoped {
	if len(log) == 0 {
		return &debugScopedLock{log: func(isAboutToLock bool) {
			if isAboutToLock {
				fmt.Println("Scoped Lock is About to be Locked")
			} else {
				fmt.Println("Scoped Lock is About to be Unlocked")
			}
		}}
	}
	return &debugScopedLock{log: log[0]}
}

// Lock locks the ScopedLock. This can be used in a 'with' statement.
func (s *scopedLock) Lock() func() {
	s.mu.Lock()
	return func() { s.mu.Unlock() }
}
