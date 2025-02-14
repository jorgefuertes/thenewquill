package item

import "errors"

var (
	ErrContainerCantCarrySoMuch = errors.New("container can't carry so much weight")
	ErrContainerIsFull          = errors.New("container is full")
	ErrNotContainer             = errors.New("item is not a container")
	ErrDuplicateLabel           = errors.New("duplicated label")
	ErrDuplicateNounAdj         = errors.New("duplicated noun/adjective")
	ErrEmptyLabel               = errors.New("empty label")
	ErrNounCannotBeNil          = errors.New("noun cannot be nil")
	ErrNounCannotBeUnderscore   = errors.New("noun cannot be underscore")
)
