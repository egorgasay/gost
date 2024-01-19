package main

import (
	"strconv"

	"github.com/egorgasay/gost"
)

const (
	UnknownBaseCode = iota
	InvalidBaseCode
)

const (
	UnknownExtCode = iota
	FormatExtCode
	ParamExtCode
)

func doSomeWork(x int) gost.Result[gost.Nothing] {
	// ...

	if x == 2 {
		return gost.Err[gost.Nothing](
			gost.NewError(InvalidBaseCode, FormatExtCode, "invalid x"),
		)
	}

	return gost.Ok(gost.Nothing{})
}

func convertUserID(userID string) gost.Result[uint] {
	userID64, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return gost.Err[uint](
			gost.NewError(InvalidBaseCode, FormatExtCode, "invalid user id"),
		)
	}

	return gost.Ok(uint(userID64))
}

func convertCoords(lat, lon string) (res gost.Result[gost.Pair[float64, float64]]) {
	latF, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return res.Err(
			gost.NewError(InvalidBaseCode, FormatExtCode, "invalid latitude"),
		)
	}

	lonF, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		return res.Err(
			gost.NewError(InvalidBaseCode, FormatExtCode, "invalid longitude"),
		)
	}

	return res.Ok(gost.Pair[float64, float64]{latF, lonF})
}

func returnValueIfNotNil(x *int) gost.Option[int] {
	if x == nil {
		return gost.None[int]()
	}

	return gost.Some(*x)
}

func div(x, y int) gost.Result[int] {
	if y == 0 {
		return gost.Err[int](
			gost.NewError(InvalidBaseCode, ParamExtCode, "zero division"),
		)
	}

	return gost.Ok(x / y)
}

func pairReturned() gost.Result[gost.Pair[int, string]] {
	return gost.Ok(gost.Pair[int, string]{1, "2"})
}

func tupleReturned() gost.Result[gost.Tuple[int]] {
	return gost.Ok(gost.Tuple[int]{1, 2})
}

func convertErr(err *gost.Error) {
	if err == nil {
		return
	}

	base := err.BaseCode()
	extended := err.ExtendedCode()

	switch base {
	case InvalidBaseCode:
		switch extended {
		case FormatExtCode:
			// do something
			// ...
		}
		// ....
	}
}

func returnOption() (opt gost.Option[int]) {
	return opt.Some(1)
}

func returnOptionNone() (opt gost.Option[int]) {
	return opt.None()
}

func returnOptionV2() gost.Option[int] {
	return gost.Some(1)
}

func returnOptionNoneV2() gost.Option[int] {
	return gost.None[int]()
}

func returnResult() (res gost.Result[int]) {
	return res.Ok(1)
}

func returnResultError() (opt gost.Result[int]) {
	return opt.Err(gost.ErrMutexIsLocked)
}

func returnResultV2() gost.Result[int] {
	return gost.Ok(1)
}

func returnResultErrorV2() gost.Result[int] {
	return gost.Err[int](gost.ErrMutexIsLocked)
}

func main() {
	//var s = gost.NewDebugScopedLock()
	//var m = map[int]struct{}{}
	//
	//for i := 0; i < 20000; i++ {
	//	go func() {
	//		doWithLock(s, m)
	//	}()
	//}
	//
	//time.Sleep(2 * time.Second)

	mu := gost.NewSecureMutex("key")

	defer mu.Lock().Unlock()
}

func doWithLock(s gost.Scoped, m map[int]struct{}) {
	defer s.Lock()() // lock here

	// Critical section
	m[1] = struct{}{}

	// unlock here
}
