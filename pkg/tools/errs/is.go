package errs

import (
	"errors"
	"net/http"
)

func Is(err error) (Error, bool) {
	var e Error
	return e, errors.As(err, &e)
}

func IsForbidden(err error) bool {
	var e Error
	return errors.As(err, &e) && e.HTTPStatusCode == http.StatusForbidden
}

func IsNotFound(err error) bool {
	var e Error
	return errors.As(err, &e) && e.HTTPStatusCode == http.StatusNotFound
}

func IsConflict(err error) bool {
	var e Error
	return errors.As(err, &e) && e.HTTPStatusCode == http.StatusConflict
}

func IsTooManyReqs(err error) bool {
	var e Error
	return errors.As(err, &e) && e.HTTPStatusCode == http.StatusTooManyRequests
}
