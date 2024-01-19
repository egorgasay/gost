package gost

import (
	"encoding/json"
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
	extCodes []int
	messages []string
	unwrap   []error
}

func (x *ErrX) BaseCode() int {
	if x == nil {
		panic("base code is nil")
	}

	return x.baseCode
}

func (x ErrX) ExtCodes() []int {
	return x.extCodes
}

func (x ErrX) Messages() []string {
	return x.messages
}

func NewErrX(code int, message ...string) *ErrX {
	err := ErrX{
		baseCode: code,
		extCodes: make([]int, 0),
		messages: message,
		unwrap:   make([]error, 0), // TODO:
	}

	if len(message) == 0 {
		message = []string{""}
	}

	return &err
}

//var Nil = ErrX{baseCode: nil}

func (x *ErrX) Extend(extCode int, messages ...string) *ErrX {
	if x == nil {
		return nil
	}

	var message string

	if len(messages) > 0 {
		message = strings.Join(messages, ", ")
	}

	x.extCodes = append(x.extCodes, extCode)
	x.messages = append(x.messages, message)

	return x
}

func (x *ErrX) Unwrap() error {
	return x
}

func (x *ErrX) ExtendMsg(message string, messages ...string) *ErrX {
	if x == nil {
		return nil
	}

	if len(messages) > 0 {
		message += fmt.Sprintf("; %s", strings.Join(messages, ", "))
	}

	x.extCodes = append(x.extCodes, 0),
		messages: append(x.messages, message),
	}
}

type extendedCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type errXJSON struct {
	BaseCode      int            `json:"base_code"`
	Message       string         `json:"message"`
	ExtendedCodes []extendedCode `json:"extended_codes"`
}

// MarshalJSON implements the json.Marshaler interface.
// Example:
//
//		{
//			"base_code": 1,
//			"message": "Not found",
//			"extended_codes": [
//				{
//					"code": 234,
//					"message": "Order"
//				},
//				{
//					"code": 0,
//					"message": "can't find order"
//				}
//			],
//	}
func (x *ErrX) MarshalJSON() ([]byte, error) {
	if x == nil {
		return []byte("null"), nil
	}

	// Create the JSON structure using the fields from ErrX.
	baseCode := x.baseCode
	message := "Unknown error"
	if len(x.messages) > 0 {
		message = x.messages[0]
		x.messages = x.messages[1:]
	}

	extendedCodes := make([]extendedCode, len(x.extCodes))
	for i, code := range x.extCodes {
		extendedMessage := "Unknown error"
		if i < len(x.messages) {
			extendedMessage = x.messages[i]
		}
		extendedCodes[i] = extendedCode{
			Code:    code,
			Message: extendedMessage,
		}
	}

	// Marshal the JSON structure.
	errX := errXJSON{
		BaseCode:      baseCode,
		Message:       message,
		ExtendedCodes: extendedCodes,
	}

	return json.Marshal(errX)
}

func (x *ErrX) Join(err *ErrX) *ErrX {
	if x == nil {
		return nil
	}

	return &ErrX{
		baseCode: x.baseCode,
		extCodes: append(x.extCodes, err.extCodes...),
		messages: append(x.messages, err.messages...),
		unwrap:   append(x.unwrap, err),
	}
}

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

	for _, v := range x.extCodes {
		if v == code {
			return true
		}
	}

	return false
}

func (x *ErrX) AsMessage(err error) *ErrX {
	if x == nil {
		return nil
	}

	x.extCodes = append(x.extCodes, 0)
	x.messages = append(x.messages, err.Error())

	return x
}

// Error
// Example:
// fmt.Println(newErr(_notFound, "not found").Extend(_order).Extend(134, "test").Error())
// OUT: 404: "not found", 1000: , 134: "test"
func (x *ErrX) Error() string {
	if x == nil {
		panic("not an error")
	}

	var message = fmt.Sprintf("%d: %s", x.baseCode, x.messages[0])

	if len(x.messages) > 0 {
		x.messages = x.messages[1:]
	}

	for i := 0; i < len(x.extCodes) || i < len(x.messages); i++ {
		message += "; "

		if i < len(x.extCodes) {
			message += fmt.Sprintf("%d: ", x.extCodes[i])
		}

		if i < len(x.messages) {
			message += x.messages[i]
		}
	}

	return message
}
