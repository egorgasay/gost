package gost

import (
	"encoding/json"
	"fmt"
)

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

func NewErrorUnknown(message string) *Error {
	return NewError(0, 0, message)
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

func (e *Error) WrapNotNilMsg(message string) *Error {
	if e != nil {
		e.message = fmt.Sprintf("%s: %s", message, e.message)
		return e
	}

	return nil
}

func (e *Error) WrapfNotNilMsg(template string, args ...any) *Error {
	if e != nil {
		e.message = fmt.Sprintf("%s: %s", fmt.Sprintf(template, args...), e.message)
		return e
	}

	return nil
}

func (e *Error) Message() string {
	return e.message
}

func (e Error) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"base_code":%d,"extended_code":%d,"message":"%s"}`, e.baseCode, e.extendedCode, e.message)), nil
}

func (e *Error) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, e); err != nil {
		return err
	}

	return nil
}

func (e *Error) IfErr(fn func(err *Error) *Error) *Error {
	if e == nil {
		return nil
	}

	return fn(e)
}

func (e *Error) IfNotErr(fn func() *Error) *Error {
	if e != nil {
		return e
	}

	return fn()
}

func (e *Error) IntoStd() error {
	if e == nil {
		return nil
	}

	return e
}
