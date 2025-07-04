package character

import "errors"

var (
	ErrOnlyOneHuman     = errors.New("only one human can be defined")
	ErrNoHuman          = errors.New("no human character defined")
	ErrEmptyDescription = errors.New("character has empty description")
	ErrCannotImport     = errors.New("cannot import character")
)
