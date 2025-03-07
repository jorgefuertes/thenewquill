package db

import "errors"

var (
	ErrEOF                = errors.New("end of file")
	ErrUnexpectedEOF      = errors.New("unexpected end of file")
	ErrShortRegister      = errors.New("short register with no fields")
	ErrInvalidSection     = errors.New("invalid section")
	ErrMaxHeaderLines     = errors.New("maximum header lines limit reached")
	ErrEOG                = errors.New("end of group")
	ErrHashMismatch       = errors.New("hash mismatch")
	ErrHeaderLimitReached = errors.New("maximum header lines limit reached")
)
