package character

import "errors"

var (
	ErrOnlyOneHuman     = errors.New("only one human can be defined")
	ErrNoHuman          = errors.New("no human character defined")
	ErrEmptyLabel       = errors.New("character has empty label")
	ErrWrongLabel       = errors.New("character has wrong label")
	ErrEmptyDescription = errors.New("character has empty description")
)
