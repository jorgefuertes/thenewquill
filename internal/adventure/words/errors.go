package words

import "errors"

var (
	ErrUnknownWordType = errors.New("unknown word type")
	ErrEmptyLabel      = errors.New("empty label")
	ErrDuplicatedWord  = errors.New("duplicated type and synonym")
)
