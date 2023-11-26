package xres_test

import (
	xres "github.com/egorgasay/xres"
	"strconv"
)

const (
	UnknownBaseCode = iota
	InvalidBaseCode = iota
)

const (
	UnknownExtCode = iota
	FormatExtCode  = iota
	ParamExtCode   = iota
)

func doSomeWork(x int) xres.Result[xres.Nothing] {
	// ...

	if x == 2 {
		return xres.Err[xres.Nothing](xres.NewError(InvalidBaseCode, FormatExtCode, "invalid x"))
	}

	return xres.Ok(xres.Nothing{})
}

func convertUserID(userID string) xres.Result[uint] {
	userID64, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return xres.Err[uint](
			xres.NewError(InvalidBaseCode, FormatExtCode, "invalid user id"),
		)
	}

	return xres.Ok(uint(userID64))
}

func convertCoords(lat, lon string) xres.Result[xres.Pair[float64, float64]] {
	latF, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return xres.Err[xres.Pair[float64, float64]](
			xres.NewError(InvalidBaseCode, FormatExtCode, "invalid latitude"),
		)
	}

	lonF, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		return xres.Err[xres.Pair[float64, float64]](
			xres.NewError(InvalidBaseCode, FormatExtCode, "invalid longitude"),
		)
	}

	return xres.Ok(xres.Pair[float64, float64]{latF, lonF})
}

func returnValueIfNotNil(x *int) xres.Option[int] {
	if x == nil {
		return xres.None[int]()
	}

	return xres.Some(*x)
}

func div(x, y int) xres.Result[int] {
	if y == 0 {
		return xres.Err[int](
			xres.NewError(InvalidBaseCode, ParamExtCode, "zero division"),
		)
	}

	return xres.Ok(x / y)
}

func pairReturned() xres.Result[xres.Pair[int, string]] {
	return xres.Ok(xres.Pair[int, string]{1, "2"})
}

func convertErr(err *xres.Error) {
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
	_ := pairReturned().UnwrapOrDefault().Left
}
