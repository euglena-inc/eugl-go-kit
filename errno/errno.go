package errno

import "errors"

const (
	CodeSuccess       = 0
	CodeInvalidParam  = 1001
	CodeUnauthorized  = 1002
	CodeForbidden     = 1003
	CodeNotFound      = 1004
	CodeConflict      = 1005
	CodeInternalError = 1006
)

type Error struct {
	Code    int
	Message string
	Cause   error
}

func New(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

func Wrap(code int, message string, cause error) *Error {
	return &Error{Code: code, Message: message, Cause: cause}
}

func (e *Error) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return e.Message + ": " + e.Cause.Error()
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func As(err error) (*Error, bool) {
	var target *Error
	ok := errors.As(err, &target)
	return target, ok
}
