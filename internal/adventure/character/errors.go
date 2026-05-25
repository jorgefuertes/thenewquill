package character

import "errors"

var (
	ErrOnlyOneHuman     = errors.New("only one human can be defined")
	ErrNoHuman          = errors.New("no human character defined")
	ErrEmptyDescription = errors.New("character has empty description")
	ErrCannotImport     = errors.New("cannot import character")
	ErrMissingNoun      = errors.New("missing noun")
	ErrMissingAdjective = errors.New("missing adjective")
	ErrItemHasSameNoun  = errors.New("there's an item with the same name")
	ErrItemHasSameLabel = errors.New("there's an item with the same label")
)
