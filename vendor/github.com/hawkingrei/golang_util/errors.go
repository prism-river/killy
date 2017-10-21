// it is an example to show the best practices for errors
package util

import (
	"fmt"
	"net/http"
	"net/textproto"
	"reflect"
	"runtime"
	"strings"
	"unicode/utf8"
)

type HTTPError interface {
	// Error returns error's message.
	Error() string
	// Status returns error's http status code.
	Status() int
}
// Error represents a numeric error with optional meta. It can be used in middleware as a return result.
type Error struct {
	Code  int         `json:"-"`
	Err   string      `json:"error"`
	Msg   string      `json:"message"`
	Data  interface{} `json:"data,omitempty"`
	Stack string      `json:"-"`
}

// Status implemented HTTPError interface.
func (err *Error) Status() int {
	return err.Code
}

// Error implemented HTTPError interface.
func (err *Error) Error() string {
	return fmt.Sprintf("%s: %s", err.Err, err.Msg)
}

// String implemented fmt.Stringer interface.
func (err Error) String() string {
	return err.GoString()
}

// GoString implemented fmt.GoStringer interface, returns a Go-syntax string.
func (err Error) GoString() string {
	if v, ok := err.Data.([]byte); ok && utf8.Valid(v) {
		err.Data = string(v)
	}
	return fmt.Sprintf(`Error{Code:%d, Err:"%s", Msg:"%s", Data:%#v, Stack:"%s"}`,
		err.Code, err.Err, err.Msg, err.Data, err.Stack)
}

// WithMsg returns a copy of err with given new messages.
//  err := gear.Err.WithMsg() // just clone
//  err := gear.ErrBadRequest.WithMsg("invalid email") // 400 Bad Request error with message invalid email"
func (err Error) WithMsg(msgs ...string) *Error {
	if len(msgs) > 0 {
		err.Msg = strings.Join(msgs, ", ")
	}
	return &err
}

// WithMsgf returns a copy of err with given message in the manner of fmt.Printf.
//  err := gear.ErrBadRequest.WithMsgf(`invalid email: "%s"`, email)
func (err Error) WithMsgf(format string, args ...interface{}) *Error {
	return err.WithMsg(fmt.Sprintf(format, args...))
}

// WithCode returns a copy of err with given code.
//  BadRequestErr := gear.Err.WithCode(400)
func (err Error) WithCode(code int) *Error {
	err.Code = code
	if text := http.StatusText(code); text != "" {
		err.Err = text
	}
	return &err
}

// WithStack returns a copy of err with error stack.
//  err := gear.Err.WithMsg("some error").WithStack()
func (err Error) WithStack(skip ...int) *Error {
	return ErrorWithStack(&err, skip...)
}

// From returns a copy of err with given error. It will try to merge the given error.
// If the given error is a *Error instance, it will be returned without copy.
//  err := gear.ErrBadRequest.From(errors.New("invalid email"))
//  err := gear.Err.From(someErr)
func (err Error) From(e error) *Error {
	if IsNil(e) {
		return nil
	}

	switch v := e.(type) {
	case *Error:
		return v
	case HTTPError:
		err.Code = v.Status()
		err.Msg = v.Error()
	case *textproto.Error:
		err.Code = v.Code
		err.Msg = v.Msg
	default:
		err.Msg = e.Error()
	}

	if err.Err == "" {
		err.Err = http.StatusText(err.Code)
	}
	return &err
}

// ParseError parse a error, textproto.Error or HTTPError to HTTPError
func ParseError(e error, code ...int) HTTPError {
	if IsNil(e) {
		return nil
	}

	switch v := e.(type) {
	case HTTPError:
		return v
	case *textproto.Error:
		err := Err.WithCode(v.Code)
		err.Msg = v.Msg
		return err
	default:
		err := ErrInternalServerError.WithMsg(e.Error())
		if len(code) > 0 && code[0] > 0 {
			err = err.WithCode(code[0])
		}
		return err
	}
}

// Predefined errors
var (
	Err = &Error{Code: http.StatusInternalServerError, Err: "Error"}

	// https://golang.org/pkg/net/http/#pkg-constants
	ErrBadRequest                    = Err.WithCode(http.StatusBadRequest)
	ErrUnauthorized                  = Err.WithCode(http.StatusUnauthorized)
	ErrPaymentRequired               = Err.WithCode(http.StatusPaymentRequired)
	ErrForbidden                     = Err.WithCode(http.StatusForbidden)
	ErrNotFound                      = Err.WithCode(http.StatusNotFound)
	ErrMethodNotAllowed              = Err.WithCode(http.StatusMethodNotAllowed)
	ErrNotAcceptable                 = Err.WithCode(http.StatusNotAcceptable)
	ErrProxyAuthRequired             = Err.WithCode(http.StatusProxyAuthRequired)
	ErrRequestTimeout                = Err.WithCode(http.StatusRequestTimeout)
	ErrConflict                      = Err.WithCode(http.StatusConflict)
	ErrGone                          = Err.WithCode(http.StatusGone)
	ErrLengthRequired                = Err.WithCode(http.StatusLengthRequired)
	ErrPreconditionFailed            = Err.WithCode(http.StatusPreconditionFailed)
	ErrRequestEntityTooLarge         = Err.WithCode(http.StatusRequestEntityTooLarge)
	ErrRequestURITooLong             = Err.WithCode(http.StatusRequestURITooLong)
	ErrUnsupportedMediaType          = Err.WithCode(http.StatusUnsupportedMediaType)
	ErrRequestedRangeNotSatisfiable  = Err.WithCode(http.StatusRequestedRangeNotSatisfiable)
	ErrExpectationFailed             = Err.WithCode(http.StatusExpectationFailed)
	ErrTeapot                        = Err.WithCode(http.StatusTeapot)
	ErrUnprocessableEntity           = Err.WithCode(http.StatusUnprocessableEntity)
	ErrLocked                        = Err.WithCode(http.StatusLocked)
	ErrFailedDependency              = Err.WithCode(http.StatusFailedDependency)
	ErrUpgradeRequired               = Err.WithCode(http.StatusUpgradeRequired)
	ErrPreconditionRequired          = Err.WithCode(http.StatusPreconditionRequired)
	ErrTooManyRequests               = Err.WithCode(http.StatusTooManyRequests)
	ErrRequestHeaderFieldsTooLarge   = Err.WithCode(http.StatusRequestHeaderFieldsTooLarge)
	ErrUnavailableForLegalReasons    = Err.WithCode(http.StatusUnavailableForLegalReasons)
	ErrInternalServerError           = Err.WithCode(http.StatusInternalServerError)
	ErrNotImplemented                = Err.WithCode(http.StatusNotImplemented)
	ErrBadGateway                    = Err.WithCode(http.StatusBadGateway)
	ErrServiceUnavailable            = Err.WithCode(http.StatusServiceUnavailable)
	ErrGatewayTimeout                = Err.WithCode(http.StatusGatewayTimeout)
	ErrHTTPVersionNotSupported       = Err.WithCode(http.StatusHTTPVersionNotSupported)
	ErrVariantAlsoNegotiates         = Err.WithCode(http.StatusVariantAlsoNegotiates)
	ErrInsufficientStorage           = Err.WithCode(http.StatusInsufficientStorage)
	ErrLoopDetected                  = Err.WithCode(http.StatusLoopDetected)
	ErrNotExtended                   = Err.WithCode(http.StatusNotExtended)
	ErrNetworkAuthenticationRequired = Err.WithCode(http.StatusNetworkAuthenticationRequired)
)

// ErrorWithStack create a error with stacktrace
func ErrorWithStack(val interface{}, skip ...int) *Error {
	if IsNil(val) {
		return nil
	}

	var err *Error
	switch v := val.(type) {
	case *Error:
		err = v.WithMsg() // must clone, should not change the origin *Error instance
	case error:
		err = ErrInternalServerError.From(v)
	case string:
		err = ErrInternalServerError.WithMsg(v)
	default:
		err = ErrInternalServerError.WithMsgf("%#v", v)
	}

	if err.Stack == "" {
		buf := make([]byte, 2048)
		buf = buf[:runtime.Stack(buf, false)]
		s := 1
		if len(skip) != 0 {
			s = skip[0]
		}
		err.Stack = pruneStack(buf, s)
	}
	return err
}

func pruneStack(stack []byte, skip int) string {
	// remove first line
	// `goroutine 1 [running]:`
	lines := strings.Split(string(stack), "\n")[1:]
	newLines := make([]string, 0, len(lines)/2)

	num := 0
	for idx, line := range lines {
		if idx%2 == 0 {
			continue
		}
		skip--
		if skip >= 0 {
			continue
		}
		num++

		loc := strings.Split(line, " ")[0]
		loc = strings.Replace(loc, "\t", "\\t", -1)
		// only need odd line
		newLines = append(newLines, loc)
		if num == 10 {
			break
		}
	}
	return strings.Join(newLines, "\\n")
}

// IsNil checks if a specified object is nil or not, without Failing.
func IsNil(val interface{}) bool {
	if val == nil {
		return true
	}

	value := reflect.ValueOf(val)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
