package db

import "errors"

var (
	ErrUnexpectedEOF      = errors.New("unexpected end of file")
	ErrShortRegister      = errors.New("short register with no fields")
	ErrInvalidSection     = errors.New("invalid section")
	ErrHashMismatch       = errors.New("hash mismatch")
	ErrHeaderLimitReached = errors.New("maximum header lines limit reached")
	ErrConfigSectionRegs  = errors.New("config section must have exactly one register")
)
