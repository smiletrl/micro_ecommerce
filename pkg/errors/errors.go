package errors

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

// Error represents business error
type Error struct {
	// this code may not be necessary now.
	// Code is for future usage, something like `invalid_username`,
	// `unmatched_password`.
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

// New is to create a new error
func New(message string, code ...string) error {
	if len(code) == 1 {
		return &Error{
			Code:    code[0],
			Message: message,
		}
	}
	// This is the default error
	return &Error{
		Code:    "error",
		Message: message,
	}
}

// Response is error response
type Response struct {
	*Error
}

// Abort means error out `500`
func Abort(c echo.Context, err error) error {
	log.Printf("%+v", errors.WithStack(err))
	// Get the cause error.
	causeErr := errors.Cause(err)

	// Assert this err type. If this is our own bussiness error, we can
	// repond accordingly.
	// Ideally we should always catch one business error from error stack.
	if rootErr, ok := causeErr.(*Error); ok {
		return c.JSON(http.StatusInternalServerError, Response{
			Error: rootErr,
		})
	}

	// We should barely reach this code finally.
	// Important, we don't show internal error to frontend directly
	return c.JSON(http.StatusInternalServerError, Response{
		Error: &Error{
			Code:    "error",
			Message: "unknown error",
		},
	})
}

// BadRequest means bad request `400`
func BadRequest(c echo.Context, err error) error {
	log.Printf("%+v", errors.WithStack(err))
	// Get the cause error.
	causeErr := errors.Cause(err)

	// Assert this err type. If this is our own bussiness error, we can
	// repond accordingly.
	// Ideally we should always catch one business error from error stack.
	if rootErr, ok := causeErr.(*Error); ok {
		return c.JSON(http.StatusBadRequest, Response{
			Error: rootErr,
		})
	}

	// We should barely reach this code finally.
	// Important, we don't show internal error to frontend directly
	return c.JSON(http.StatusBadRequest, Response{
		Error: &Error{
			Code:    "error",
			Message: err.Error(),
		},
	})
}
