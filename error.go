package gost

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Error struct {
	baseCode     int
	extendedCode int
	message      string

	joined []*Error
}

func NewError(baseCode int, extendedCode int, message string) *Error {
	return &Error{
		baseCode:     baseCode,
		extendedCode: extendedCode,
		message:      message,
		joined:       make([]*Error, 0),
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

func (e *Error) Wrap(message string) *Error {
	if e != nil {
		e.message = fmt.Sprintf("%s: %s", message, e.message)
		return e
	}

	return nil
}

func (e *Error) Wrapf(template string, args ...any) *Error {
	if e != nil {
		e.message = fmt.Sprintf("%s: %s", fmt.Sprintf(template, args...), e.message)
		return e
	}

	return nil
}

func (e *Error) Message() string {
	return e.message
}

func (e Error) MarshalJSON() ([]byte, error) { // TODO: inner errors
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

func (e *Error) Joined() []*Error {
	return e.joined
}

func (e *Error) Join(err *Error) *Error {
	if e == nil {
		return err
	}

	if err == nil {
		return e
	}

	e.joined = append(e.joined, err)

	return e
}

func (e *Error) Unwrap() []*Error {
	if e == nil {
		return nil
	}

	return append(e.joined[:0:0], e.joined...)
}

func (e *Error) Is(target *Error) bool {
	if e == target {
		return true
	}

	if e == nil || target == nil {
		return e == target
	}

	found := false

	for _, err := range e.Unwrap() {
		if err == nil {
			continue
		}

		if err == target {
			found = true
			break
		}
	}

	return found
}

type ErrX struct {
	baseCode int
	extCode  int
	message  string
	parent   *ErrX
}

const (
	ErrXSeparator        = ";;"
	ErrXMessageSeparator = " : "
)

func NewErrX(code int, message ...string) *ErrX {
	err := ErrX{
		baseCode: code,
	}

	if len(message) > 0 {
		err.message = strings.Join(message, ErrXMessageSeparator)
	}

	return &err
}

func (x *ErrX) BaseCode() int {
	if x == nil {
		panic("base code is nil")
	}

	return x.baseCode
}

func (x ErrX) ExtCode() int {
	return x.extCode
}

func (x ErrX) Message() string {
	return x.message
}

func (x ErrX) Messages() []string {
	if x.parent == nil {
		return []string{x.message}
	}

	return append(x.parent.Messages(), x.message)
}

func (x ErrX) MessagesSpace() string {
	return strings.TrimSpace(strings.Join(x.Messages(), " "))
}

//var Nil = ErrX{baseCode: nil}

func (x *ErrX) Extend(extCode int, messages ...string) *ErrX {
	if x == nil {
		return nil
	}

	err := ErrX{
		baseCode: x.baseCode,
		extCode:  extCode,
		parent:   x,
	}

	if len(messages) == 1 {
		err.message = messages[0]
	} else if len(messages) > 0 {
		err.message = strings.Join(messages, ErrXMessageSeparator)
	}

	return &err
}

func (x *ErrX) Is(target error) bool {
	if x == nil {
		return false
	}

	var errX *ErrX
	if !errors.As(target, &errX) {
		return false
	}

	return x.IsX(errX)
}

func (x *ErrX) IsX(targetX *ErrX) bool {
	if targetX == nil {
		return x == nil
	}

	if x.baseCode != targetX.baseCode {
		return false
	}

	if x.extCode != targetX.extCode {
		if x.parent == nil {
			for targetX.parent != nil {
				if x == targetX.parent {
					return true
				}

				targetX = targetX.parent
			}
		}

		return x.parent.Is(targetX)
	}

	return true
}

func (x *ErrX) ExtendMsg(message string, messages ...string) *ErrX {
	if x == nil {
		return nil
	}

	if len(messages) > 0 {
		message = fmt.Sprintf("%s %s%s%s", message, ErrXMessageSeparator, strings.Join(messages, ErrXMessageSeparator), ErrXSeparator)
	}

	return x.Extend(0, message)
}

type errXJSON struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Parent  *ErrX  `json:"parent,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface.
// Example:
//
//	{
//		"code": 12,
//		"message": "Order",
//		"parent": [
//			{
//				"code": 1,
//				"message": "Not found",
//			}
//		],
//	}
func (x *ErrX) MarshalJSON() ([]byte, error) {
	if x == nil {
		return []byte("null"), nil
	}

	errX := errXJSON{
		Code:    x.baseCode,
		Message: x.message,
		Parent:  x.parent,
	}

	if x.extCode != 0 || x.parent != nil {
		errX.Code = x.extCode
	}

	return json.Marshal(errX)
}

//func (x *ErrX) Join(err *ErrX) *ErrX {
//	if x == nil {
//		return nil
//	}
//
//	x.extCodes = append(x.extCodes, 0)
//	x.messages = append(x.messages, err.messages...)
//	x.unwrap = append(x.unwrap, err.unwrap...)
//
//	return x
//}

func (x *ErrX) CmpBase(code int) bool {
	if x == nil {
		return false
	}

	return x.baseCode == code
}

func (x *ErrX) CmpExt(code int) bool {
	if x == nil {
		return false
	}

	return x.extCode == code
}

func (x *ErrX) HasExt(code int) bool {
	if x == nil {
		return false
	}

	if x.extCode != code {
		if x.parent == nil {
			return false
		}

		return x.parent.HasExt(code)
	}

	return x.extCode == code
}

func (x *ErrX) AsMessage(err error) *ErrX {
	if err == nil {
		return nil
	}

	if x == nil {
		return NewErrX(0, err.Error())
	}

	return x.Extend(0, err.Error())
}

// Error
// Example:
// fmt.Println(newErr(_notFound, "not found").Extend(_order).Extend(134, "test").Error())
// OUT: 404: "not found";; 1000: ;; 134: "test"
func (x *ErrX) Error() string {
	if x == nil {
		panic("not an error")
	}

	code := x.baseCode
	if x.extCode != 0 {
		code = x.extCode
	}

	message := fmt.Sprintf("%d: %s", code, x.message)
	if x.parent == nil {
		return message
	}

	return fmt.Sprintf("%s%s%s", x.parent.Error(), ErrXSeparator, message)
}
