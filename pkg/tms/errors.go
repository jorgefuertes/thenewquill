package tms

import "errors"

var (
	ErrBufferTooShort = errors.New("buffer is too short")
	ErrInvalidKey     = errors.New("invalid key")
)
