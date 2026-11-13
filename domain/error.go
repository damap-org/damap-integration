package domain

import "errors"

// ErrNotFound indicates that the requested resource was not found.
var ErrNotFound = errors.New("not found")

func (e Error) Error() string {
	return e.Message
}
