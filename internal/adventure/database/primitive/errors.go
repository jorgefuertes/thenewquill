package primitive

import "errors"

var (
	ErrUndefinedID      = errors.New("undefined ID")
	ErrInvalidID        = errors.New("invalid ID")
	ErrInvalidLabelName = errors.New("invalid label name")
	ErrInvalidLabelType = errors.New("invalid label type")
)
