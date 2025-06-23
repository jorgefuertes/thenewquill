package word

import "errors"

var (
	ErrDuplicatedWord = errors.New("duplicated type and synonym")
	ErrEmptyWord      = errors.New("empty word")
)
