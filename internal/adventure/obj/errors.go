package obj

import "errors"

var ErrContainerCantCarrySoMuch = errors.New("container can't carry so much weight")
var ErrContainerIsFull = errors.New("container is full")
var ErrNotContainer = errors.New("item is not a container")
var ErrDuplicateLabel = errors.New("duplicated label")
var ErrDuplicateNounAdj = errors.New("duplicated noun/adjective")
var ErrEmptyLabel = errors.New("empty label")
var ErrNounCannotBeNil = errors.New("noun cannot be nil")
