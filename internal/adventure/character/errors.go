package character

import "errors"

var (
	ErrDuplicatedPlayerLabel = errors.New("duplicated player label")
	ErrOnlyOneHuman          = errors.New("only one human can be defined")
)
