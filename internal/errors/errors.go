package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
	Err        error  `json:"-"`
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(code int, message string) *Error {
	return &Error{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatusFromCode(code),
	}
}

func Wrap(err error, code int, message string) *Error {
	return &Error{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatusFromCode(code),
		Err:        err,
	}
}

func Is(err error, target *Error) bool {
	var appErr *Error
	if !errors.As(err, &appErr) {
		return false
	}
	return appErr.Code == target.Code
}

func httpStatusFromCode(code int) int {
	switch code {
	case CodeBadRequest:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

const (
	CodeOK            = 0
	CodeBadRequest    = 40000
	CodeUnauthorized  = 40100
	CodeForbidden     = 40300
	CodeNotFound      = 40400
	CodeConflict      = 40900
	CodeInternalError = 50000
)

var (
	ErrBadRequest    = New(CodeBadRequest, "bad request")
	ErrUnauthorized  = New(CodeUnauthorized, "unauthorized")
	ErrNotFound      = New(CodeNotFound, "resource not found")
	ErrConflict      = New(CodeConflict, "resource conflict")
	ErrInternalError = New(CodeInternalError, "internal server error")
)
