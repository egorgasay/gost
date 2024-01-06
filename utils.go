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

	var (
		once sync.Once
		done = func() { close(ch) }
	)

	if len(onStop) > 0 {
		defer onStop[0]()
	}

	pool <- struct{}{}
	go func() {
		defer func() {
			<-pool // Ensure that we release the resource in the pool
		}()

		err = fn()
		once.Do(done)
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

// TODO: use semaphore instead of channel
// TODO: GPT: Однако стоит отметить, что в вашем коде есть потенциальная проблема с обработкой ошибок. Если функция onStop вернет ошибку, эта ошибка будет перезаписана значением errFunc, если errFunc не равна nil. Это может привести к тому, что ошибка, возвращаемая функцией onStop, будет проигнорирована. Возможно, вам стоит пересмотреть эту часть кода, чтобы обе ошибки были должным образом обработаны.
//func WithContextPoolV2(ctx context.Context, fn func() error, pool chan struct{}, onStop func() error) (err error) {
//	ch := make(chan struct{})
//
//	var errFnLock = NewRwLock((error)(nil))
//
//	var (
//		once sync.Once
//		done = func() { close(ch) }
//	)
//
//	if onStop != nil {
//		defer func() {
//			errStop := onStop()
//			if err != nil {
//				err = errors.Join(err, errStop)
//			} else {
//				err = errStop
//			}
//		}()
//	}
//
//	pool <- struct{}{}
//	go func(errFnLock *RwLock[error]) {
//		defer func() {
//			<-pool // Ensure that we release the resource in the pool
//		}()
//
//		errFunc := fn()
//		errFnLock.SetWithLock(errFunc)
//		once.Do(done)
//	}(&errFnLock)
//
//	select {
//	case <-ch:
//		return errFnLock.v
//	case <-ctx.Done():
//		once.Do(done)
//		return ctx.Err()
//	}
//}
//
//func WithContextV2(ctx context.Context, fn func() error, onStop ...func()) (err error) {
//	ch := make(chan struct{}, 1)
//	defer close(ch)
//	return WithContextPool(ctx, fn, ch, onStop...)
//}
