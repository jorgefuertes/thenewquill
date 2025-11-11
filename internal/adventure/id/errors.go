package id

import "errors"

var (
	ErrUndefined = errors.New("undefined ID")
	ErrInvalid   = errors.New("invalid ID")
)
