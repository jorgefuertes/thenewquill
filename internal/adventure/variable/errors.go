package variable

import "errors"

var (
	ErrNilValue      = errors.New("nil value")
	ErrInvalidParent = errors.New("invalid parent")
)
