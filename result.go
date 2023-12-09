package gost

type Result[V any] struct {
	err   *Error
	value *V
}

type Nothing struct{}

type Pair[L any, R any] struct {
	Left  L
	Right R
}

type Tuple[V any] []V

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

func (Result[V]) Ok(value V) Result[V] {
	return Result[V]{
		value: &value,
	}
}

func (Result[V]) Err(err *Error) Result[V] {
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

func (r Result[V]) Error() *Error {
	return r.err
}

func (r Result[V]) IsOk() bool {
	return r.err == nil
}

func (r Result[V]) IsErr() bool {
	return r.err != nil
}

type switchOption bool

const (
	IsOk  switchOption = true
	IsErr switchOption = false
)

func (r Result[V]) Switch() switchOption {
	return r.err != nil
}

//func (r Result[V]) Separate() (v V, err *Error) {
//	if r.value == nil {
//		return v, r.err
//	}
//
//	return *r.value, r.err
//}
