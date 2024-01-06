package gost

import (
	"context"
	"errors"
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

func WithContextPool(ctx context.Context, fn func() error, pool chan struct{}, onStop ...func() error) (err error) {
	ch := make(chan error, 1)

	if len(onStop) > 0 {
		defer func() {
			errStop := onStop[0]()
			if errStop != nil {
				err = errors.Join(err, errStop)
			}
		}()
	}

	pool <- struct{}{}
	go func() {
		defer func() {
			close(ch)
			<-pool // Ensure that we release the resource in the pool
		}()

		ch <- fn()
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func WithContext(ctx context.Context, fn func() error, onStop ...func() error) (err error) {
	ch := make(chan struct{}, 1)
	defer close(ch)
	return WithContextPool(ctx, fn, ch, onStop...)
}
