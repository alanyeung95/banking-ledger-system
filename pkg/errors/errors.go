package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode using int to respresent the error
type ErrorCode int

// A list of common expected errors.
const (
	BadRequest       ErrorCode = 400
	Unauthorized     ErrorCode = 401
	ResourceNotFound ErrorCode = 404
	ServerError      ErrorCode = 500
)

// Error specifies the interfaces required by an error in the system.
type Error interface {
	Code() ErrorCode
	Message() string
	error
}

// genericError is an implementation of Error that contains
// an code and error message.
type genericError struct {
	code    ErrorCode
	message string
}

func (e *genericError) Code() ErrorCode {
	return e.code
}

func (e *genericError) Message() string {
	return e.message
}

func (e *genericError) Error() string {
	return fmt.Sprintf("%v : %v", e.code, e.message)
}

func (e *genericError) StatusCode() int {
	// http status code reference: https://golang.org/pkg/net/http/
	httpStatus, ok := map[ErrorCode]int{
		ResourceNotFound: http.StatusNotFound,
		BadRequest:       http.StatusBadRequest,
		Unauthorized:     http.StatusUnauthorized,
	}[e.Code()]
	if !ok {
		httpStatus = http.StatusInternalServerError
	}
	return httpStatus
}

func NewResourceNotFoundError(err error) Error {
	return &genericError{
		code:    ResourceNotFound,
		message: err.Error(),
	}
}

func NewServerError(err error) Error {
	return &genericError{
		code:    ServerError,
		message: err.Error(),
	}
}

func NewBadRequestError(err error) Error {
	return &genericError{
		code:    BadRequest,
		message: err.Error(),
	}
}

func NewUnauthorizedError(err error) Error {
	return &genericError{
		code:    Unauthorized,
		message: err.Error(),
	}
}
