package gost

import (
	"context"
	"sync"
)

func CloneArray[S ~[]E, E any](s S) S {
	return append(s[:0:0], s...)
}

func SafeDeref[T any](t *T) T {
	if t == nil {
		return *new(T)
	}
	return *t
}

func WithContextPool(ctx context.Context, fn func() error, pool chan struct{}, onStop ...func()) (err error) {
	ch := make(chan struct{})

	if len(onStop) > 0 {
		defer onStop[0]()
	}

	once := sync.Once{}
	done := func() { close(ch) }

	pool <- struct{}{}
	go func() {
		err = fn()
		once.Do(done)
		<-pool
	}()

	select {
	case <-ch:
		return err
	case <-ctx.Done():
		once.Do(done)
		return ctx.Err()
	}
}

func WithContext(ctx context.Context, fn func() error, onStop ...func()) (err error) {
	ch := make(chan struct{}, 1)
	defer close(ch)
	return WithContextPool(ctx, fn, ch, onStop...)
}
