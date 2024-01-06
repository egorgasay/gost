package gost

import (
	"sync"
)

type Mutex[V any] struct {
	mu *sync.Mutex
	v  V
}

type RwLock[V any] struct {
	mu *sync.RWMutex
	v  V
}

type SecureMutexLocked[V any] struct {
	mu *sync.Mutex
	v  V
}

type SecureMutexUnlocked[V any] struct {
	mu *sync.Mutex
	v  V
}

func (m SecureMutexUnlocked[V]) Lock() SecureMutexLocked[V] {
	m.mu.Lock()
	return SecureMutexLocked[V]{
		mu: m.mu,
		v:  m.v,
	}
}

func (m SecureMutexUnlocked[V]) Set(v V) SecureMutexUnlocked[V] {
	m.v = v
	return m
}

func (m SecureMutexUnlocked[V]) Value() Arc[*V, V] {
	return NewArc(&m.v)
}

func (m SecureMutexLocked[V]) Unlock() SecureMutexUnlocked[V] {
	m.mu.Unlock()
	return SecureMutexUnlocked[V]{
		mu: m.mu,
		v:  m.v,
	}
}

func NewSecureMutex[V any]() SecureMutexUnlocked[V] {
	return SecureMutexUnlocked[V]{
		mu: &sync.Mutex{},
	}
}

func NewSecureMutexV[V any](v V) SecureMutexUnlocked[V] {
	return SecureMutexUnlocked[V]{
		mu: &sync.Mutex{},
		v:  v,
	}
}

var (
	ErrMutexIsLocked = NewError(0, 0, "mutex is locked")
)

func NewMutex[V any](v V) Mutex[V] {
	return Mutex[V]{
		v:  v,
		mu: &sync.Mutex{},
	}
}

func NewRwLock[V any](v V) RwLock[V] {
	return RwLock[V]{
		v:  v,
		mu: &sync.RWMutex{},
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

func (m *Mutex[V]) UpdateWithLock(fn func(v *V)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	fn(&m.v)
}

func (m *Mutex[V]) SetWithLock(fn func(v V) V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.v = fn(m.v)
}

func (m *RwLock[V]) SetWithLock(v V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.v = v
}

func (m *RwLock[V]) WithLock(fn func(v V) V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.v = fn(m.v)
}

func (m *RwLock[V]) RBorrow() Arc[*V, V] {
	m.mu.RLock()
	return NewArc(&m.v)
}

func (m *RwLock[V]) ReadAndRelease() V {
	v := m.v
	m.mu.RUnlock()
	return v
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

func (m *RwLock[V]) WRelease() {
	m.mu.Unlock()
}
