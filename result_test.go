package rusty_test

import (
	"github.com/egorgasay/rusty"
	"strconv"
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

func doSomeWork(x int) rusty.Result[rusty.Nothing] {
	// ...

	if x == 2 {
		return rusty.Err[rusty.Nothing](
			rusty.NewError(InvalidBaseCode, FormatExtCode, "invalid x"),
		)
	}

	return rusty.Ok(rusty.Nothing{})
}

func convertUserID(userID string) rusty.Result[uint] {
	userID64, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return rusty.Err[uint](
			rusty.NewError(InvalidBaseCode, FormatExtCode, "invalid user id"),
		)
	}

	return rusty.Ok(uint(userID64))
}

func convertCoords(lat, lon string) rusty.Result[rusty.Pair[float64, float64]] {
	latF, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return rusty.Err[rusty.Pair[float64, float64]](
			rusty.NewError(InvalidBaseCode, FormatExtCode, "invalid latitude"),
		)
	}

	lonF, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		return rusty.Err[rusty.Pair[float64, float64]](
			rusty.NewError(InvalidBaseCode, FormatExtCode, "invalid longitude"),
		)
	}

	return rusty.Ok(rusty.Pair[float64, float64]{latF, lonF})
}

func returnValueIfNotNil(x *int) rusty.Option[int] {
	if x == nil {
		return rusty.None[int]()
	}

	return rusty.Some(*x)
}

func div(x, y int) rusty.Result[int] {
	if y == 0 {
		return rusty.Err[int](
			rusty.NewError(InvalidBaseCode, ParamExtCode, "zero division"),
		)
	}

	return rusty.Ok(x / y)
}

func pairReturned() rusty.Result[rusty.Pair[int, string]] {
	return rusty.Ok(rusty.Pair[int, string]{1, "2"})
}

func tupleReturned() rusty.Result[rusty.Tuple[int]] {
	return rusty.Ok(rusty.Tuple[int]{1, 2})
}

func convertErr(err *rusty.Error) {
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

func x() {
	_ = pairReturned().UnwrapOrDefault().Left

	_ = tupleReturned().Unwrap()[1]

	x := rusty.NewRwLock(rusty.Error{})

	{
		x.RBorrow().Read().MarshalJSON()
	}
}
