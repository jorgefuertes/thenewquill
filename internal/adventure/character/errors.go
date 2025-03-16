package character

import "errors"

var (
	ErrDuplicatedPlayerLabel = errors.New("duplicated player label")
	ErrOnlyOneHuman          = errors.New("only one human can be defined")
	ErrNoHuman               = errors.New("no human character defined")
	ErrAlreadyLocked         = errors.New("character store is already locked")
)
