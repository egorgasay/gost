package gost

import "sync"

type MutexSlice[V any] struct {
	mu *sync.RWMutex
	s  []V
}

func NewMutexSlice[V any](s []V) *MutexSlice[V] {
	return &MutexSlice[V]{
		mu: &sync.RWMutex{},
		s:  s,
	}
}

func (s *MutexSlice[V]) RPush(v V) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s = append(s.s, v)
	return len(s.s)
}

func (s *MutexSlice[V]) LPush(v V) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s = append([]V{v}, s.s...)
	return len(s.s)
}

func (s *MutexSlice[V]) RPop() V {
	s.mu.Lock()
	defer s.mu.Unlock()
	v := s.s[len(s.s)-1]
	s.s = s.s[:len(s.s)-1]
	return v
}

func (s *MutexSlice[V]) LPop() V {
	s.mu.Lock()
	defer s.mu.Unlock()
	v := s.s[0]
	s.s = s.s[1:]
	return v
}

func (s *MutexSlice[V]) InsertAt(i int, v V) *MutexSlice[V] {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s = append(s.s, v)
	copy(s.s[i+1:], s.s[i:])
	s.s[i] = v
	return s
}

func (s *MutexSlice[V]) RemoveAt(i int) V {
	s.mu.Lock()
	defer s.mu.Unlock()
	v := s.s[i]
	s.s = append(s.s[:i], s.s[i+1:]...)
	return v
}

func (s *MutexSlice[V]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.s)
}

func (s *MutexSlice[V]) Get(i int) V {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.s[i]
}

func (s *MutexSlice[V]) Set(i int, v V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s[i] = v
}

func (s *MutexSlice[V]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s = nil
}
