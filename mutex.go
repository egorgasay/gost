package gost

import (
	"sync"
)

type Mutex[V any] struct {
	mu sync.Mutex
	v  V
}

type RwLock[V any] struct {
	mu sync.RWMutex
	v  V
}

var (
	ErrMutexIsLocked = NewError(0, 0, "mutex is locked")
)

func NewMutex[V any](v V) Mutex[V] {
	return Mutex[V]{
		v:  v,
		mu: sync.Mutex{},
	}
}

func NewRwLock[V any](v V) RwLock[V] {
	return RwLock[V]{
		v:  v,
		mu: sync.RWMutex{},
	}
}

func (m *Mutex[V]) Borrow() Arc[*V, V] {
	m.mu.Lock()
	return NewArc(&m.v)
}

func (m *Mutex[V]) BorrowMut() *V {
	m.mu.Lock()
	return &m.v
}

func (m *Mutex[V]) TryBorrow() Result[V] {
	if m.mu.TryLock() {
		return Ok(m.v)
	}

	return Err[V](ErrMutexIsLocked)
}

func (m *Mutex[V]) Return(v V) {
	m.v = v
	m.mu.Unlock()
}

func (m *Mutex[V]) Release() {
	m.mu.Unlock()
}

func (m *RwLock[V]) RBorrow() Arc[*V, V] {
	m.mu.RLock()
	return NewArc(&m.v)
}

func (m *RwLock[V]) RTryBorrow() Result[V] {
	if m.mu.TryRLock() {
		return Ok(m.v)
	}

	return Err[V](ErrMutexIsLocked)
}

func (m *RwLock[V]) Release() {
	m.mu.RUnlock()
}

func (m *RwLock[V]) WBorrow() Arc[*V, V] {
	m.mu.Lock()
	return NewArc(&m.v)
}

func (m *RwLock[V]) WTryBorrow() Result[V] {
	if m.mu.TryLock() {
		return Ok(m.v)
	}

	return Err[V](ErrMutexIsLocked)
}

func (m *RwLock[V]) WReturn(v V) {
	m.v = v
	m.mu.Unlock()
}
