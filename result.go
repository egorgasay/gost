package xres

type Result[V any] struct {
	err   *Error
	value *V
}

type Nothing struct{}

type Pair[L any, R any] struct {
	Left  L
	Right R
}

func Ok[V any](value V) Result[V] {
	return Result[V]{
		value: &value,
	}
}

func Err[V any](err *Error) Result[V] {
	return Result[V]{
		err: err,
	}
}

func (r Result[V]) Unwrap() V {
	if r.err != nil {
		panic(r.err)
	}

	return *r.value
}

func (r *Result[V]) Expect(msg string) V {
	if r.err != nil {
		panic(msg)
	}

	return *r.value
}

func (r Result[V]) UnwrapOrElse(fn func(err Error) V) V {
	if r.err != nil {
		return fn(*r.err)
	}

	return *r.value
}

func (r Result[V]) UnwrapOrDefault() V {
	return *r.value
}

func (r Result[V]) Err() *Error {
	return r.err
}
