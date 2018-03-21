package httprouter

import (
	"errors"
)

type HttpError struct {
	code int
	Err  error
}

func NewHE(code int, msg string) *HttpError {
	err := new(HttpError)
	err.code = code
	err.Err = errors.New(msg)

	return err
}

func (he *HttpError) Error() string {
	return he.Err.Error()
}

func (he *HttpError) Code() int {
	return he.code
}
