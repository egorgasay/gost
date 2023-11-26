package xres

import "fmt"

type Error struct {
	baseCode     int
	extendedCode int
	message      string
}

func NewError(baseCode int, extendedCode int, message string) *Error {
	return &Error{
		baseCode:     baseCode,
		extendedCode: extendedCode,
		message:      message,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d %d %s", e.baseCode, e.extendedCode, e.message)
}

func (e *Error) BaseCode() int {
	return e.baseCode
}

func (e *Error) ExtendedCode() int {
	return e.extendedCode
}

func (e *Error) Message() string {
	return e.message
}
