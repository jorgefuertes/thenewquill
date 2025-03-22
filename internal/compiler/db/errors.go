package db

import "errors"

var (
	ErrUnexpectedEOF       = errors.New("unexpected EOF")
	ErrUnexpectedReadError = errors.New("unexpected read error")
	ErrExpectedEOR         = errors.New("expected EOR")
	ErrHashMismatch        = errors.New("hash mismatch")
)
