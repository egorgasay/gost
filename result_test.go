package xres_test

import xres "github.com/egorgasay/xres"

const (
	UnknownBaseCode = iota
	InvalidBaseCode = iota
)

const (
	UnknownExtCode = iota
	FormatExtCode  = iota
	ParamExtCode   = iota
)

func doSomeWork(x int) xres.Result[xres.None] {
	// ...

	if x == 2 {
		return xres.Err[xres.None](xres.NewError(InvalidBaseCode, FormatExtCode, "invalid x"))
	}

	return xres.Ok(xres.None{})
}

func div(x, y int) xres.Result[int] {
	if y == 0 {
		return xres.Err[int](
			xres.NewError(InvalidBaseCode, ParamExtCode, "zero division"),
		)
	}

	return xres.Ok(x / y)
}

func multTypesReturned() xres.Result[xres.Pair[int, string]] {
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
		}
	}
}
