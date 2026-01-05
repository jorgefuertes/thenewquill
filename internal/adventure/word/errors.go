package word

import "errors"

var (
	ErrDuplicatedLabel = errors.New("duplicated label for type")
	ErrDuplicatedSyn   = errors.New("duplicated synonym for type")
	ErrWordNotFound    = errors.New("word not found")
)
