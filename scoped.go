package gost

import (
	"fmt"
	"runtime"
	"sync"
)

type scopedLock struct {
	mu sync.Mutex
}

type debugScopedLock struct {
	mu  sync.Mutex
	log func()
}

func (d *debugScopedLock) Lock() func() {
	if !d.mu.TryLock() {
		d.log()
		d.mu.Lock()
	}

	return func() {
		d.mu.Unlock()
	}
}

type Scoped interface {
	Lock() func()
}

func NewScopedLock() Scoped {
	return &scopedLock{}
}

func NewDebugScopedLock(log ...func()) Scoped {
	if len(log) == 0 {
		return &debugScopedLock{log: func() {
			_, name, line, _ := runtime.Caller(2)
			fmt.Printf("Scoped Lock is about to be Locked TWICE, in %s at line %d\n", name, line)
		}}
	}
	return &debugScopedLock{log: log[0]}
}

// Lock locks the ScopedLock. This can be used in a 'with' statement.
func (s *scopedLock) Lock() func() {
	s.mu.Lock()
	return func() { s.mu.Unlock() }
}
