package gost

import "strings"

type Result[V any] struct {
	err   *ErrX
	value *V
}

type ResultN struct {
	err *ErrX
}

type Nothing struct{}

type Pair[L any, R any] struct {
	Left  L
	Right R
}

func NewPair[L any, R any](l L, r R) Pair[L, R] {
	return Pair[L, R]{Left: l, Right: r}
}

type Tuple[V any] []V

func Ok[V any](value V) Result[V] {
	return Result[V]{
		value: &value,
	}
}

func Err[V any](err *ErrX) Result[V] {
	return Result[V]{
		err: err,
	}
}

func ErrN(err *ErrX) ResultN {
	return ResultN{err: err}
}

func (Result[V]) Ok(value V) Result[V] {
	return Result[V]{
		value: &value,
	}
}

func (ResultN) Ok() ResultN {
	return ResultN{}
}

func (Result[V]) Err(err *ErrX) Result[V] {
	return Result[V]{
		err: err,
	}
}

func (ResultN) Err(err *ErrX) ResultN {
	return ResultN{
		err: err,
	}
}

func (Result[V]) ErrNew(code, extCode int, msg ...string) Result[V] {
	fMsg, sMsg := "", ""

	if len(msg) > 0 {
		fMsg = msg[0]
		if len(msg) > 1 {
			sMsg = msg[1]

			if len(msg) > 2 {
				sMsg = strings.Join(msg[1:], ErrXMessageSeparator)
			}
		}
	}

	return Result[V]{
		err: NewErrX(code, fMsg).Extend(extCode, sMsg),
	}
}

func (ResultN) ErrNew(code, extCode int, msg ...string) ResultN {
	fMsg, sMsg := "", ""

	if len(msg) > 0 {
		fMsg = msg[0]
		if len(msg) > 1 {
			sMsg = msg[1]

			if len(msg) > 2 {
				sMsg = strings.Join(msg[1:], ErrXMessageSeparator)
			}
		}
	}

	return ResultN{
		err: NewErrX(code, fMsg).Extend(extCode, sMsg),
	}
}

func (Result[V]) ErrNewUnknown(msg string) Result[V] {
	return Result[V]{
		err: NewErrX(0, msg),
	}
}

func (ResultN) ErrNewUnknown(msg string) ResultN {
	return ResultN{
		err: NewErrX(0, msg),
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

func (r Result[V]) UnwrapOrElse(fn func(err *ErrX) V) V {
	if r.err != nil {
		return fn(r.err)
	}

	return *r.value
}

func (r Result[V]) UnwrapOrDefault() (v V) {
	if r.err != nil {
		return v
	}

	return *r.value
}

func (r Result[V]) UnwrapOr(v V) V {
	if r.err != nil {
		return v
	}

	return *r.value
}

func (r Result[V]) Error() *ErrX {
	return r.err
}

func (r Result[V]) IsOk() bool {
	return r.err == nil
}

func (r Result[V]) IsErr() bool {
	return r.err != nil
}

func (r ResultN) Error() *ErrX {
	return r.err
}

func (r ResultN) IsOk() bool {
	return r.err == nil
}

func (r ResultN) IsErr() bool {
	return r.err != nil
}

type switchOption bool

const (
	IsOk  switchOption = true
	IsErr switchOption = false
)

func (r Result[V]) Switch() switchOption {
	return r.err == nil
}

func (r ResultN) Switch() switchOption {
	return r.err == nil
}

func (r Result[V]) ErrorStd() error {
	if r.err == nil {
		return nil
	}

	return r.err
}

func (r ResultN) ErrorStd() error {
	if r.err == nil {
		return nil
	}

	return r.err
}

//func (r Result[V]) Separate() (v V, err *ErrX) {
//	if r.value == nil {
//		return v, r.err
//	}
//
//	return *r.value, r.err
//}
