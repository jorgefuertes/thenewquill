package item

import "errors"

var (
	ErrContainerCantCarrySoMuch               = errors.New("container can't carry so much weight")
	ErrContainerIsFull                        = errors.New("container is full")
	ErrNotContainer                           = errors.New("item is not a container")
	ErrDuplicateLabel                         = errors.New("duplicated label")
	ErrDuplicateNounAdj                       = errors.New("duplicated noun/adjective")
	ErrEmptyLabel                             = errors.New("empty label")
	ErrNounCannotBeNil                        = errors.New("noun cannot be nil")
	ErrAdjectiveCannotBeNil                   = errors.New("adjective cannot be nil")
	ErrNounCannotBeUnderscore                 = errors.New("noun cannot be underscore")
	ErrItemCannotBeHeldAndHaveLocation        = errors.New("item cannot be held and have location")
	ErrItemCannotBeHeldAndWorn                = errors.New("item cannot be held and worn")
	ErrItemCannotBeContainedInAndHaveLocation = errors.New("item cannot be contained in and have location")
	ErrWeightShouldBeLessOrEqualThanMaxWeight = errors.New("weight should be less or equal than max weight")
	ErrWeightCannotBeNegative                 = errors.New("weight nor max weight cannot be negative")
)
