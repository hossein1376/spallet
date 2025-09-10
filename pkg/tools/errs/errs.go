package errs

import (
	"fmt"
	"net/http"
)

type Error struct {
	Err            error
	HTTPStatusCode int
	Message        string
}

func NewErr(httpStatusCode int, opts []Options) Error {
	e := Error{HTTPStatusCode: httpStatusCode}
	for _, opt := range opts {
		opt(&e)
	}
	return e
}

func (e Error) Error() string {
	var text string
	if e.Err != nil {
		text = e.Err.Error()
	} else {
		text = http.StatusText(e.HTTPStatusCode)
	}
	return fmt.Sprintf("[%d] %s", e.HTTPStatusCode, text)
}

func (e Error) Unwrap() error {
	return e.Err
}

func BadRequest(opts ...Options) Error {
	return NewErr(http.StatusBadRequest, opts)

}

func Unauthorized(opts ...Options) Error {
	return NewErr(http.StatusUnauthorized, opts)
}

func Forbidden(opts ...Options) Error {
	return NewErr(http.StatusForbidden, opts)
}

func NotFound(opts ...Options) Error {
	return NewErr(http.StatusNotFound, opts)
}

func Conflict(opts ...Options) Error {
	return NewErr(http.StatusConflict, opts)
}

func TooMany(opts ...Options) Error {
	return NewErr(http.StatusTooManyRequests, opts)
}

func Internal(opts ...Options) Error {
	return NewErr(http.StatusInternalServerError, opts)
}

func Timeout(opts ...Options) Error {
	return NewErr(http.StatusGatewayTimeout, opts)
}
