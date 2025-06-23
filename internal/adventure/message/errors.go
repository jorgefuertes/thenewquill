package message

import "errors"

var (
	ErrUndefinedText   = errors.New("undefined text")
	ErrUndefinedPlural = errors.New("undefined plural")
)
