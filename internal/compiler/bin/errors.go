package bin

import "errors"

var (
	ErrEncodingError = errors.New("encoding error")
	ErrDecodingError = errors.New("decoding error")
)
