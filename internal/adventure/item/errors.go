package item

import "errors"

var (
	ErrContainerCantCarrySoMuch               = errors.New("container can't carry so much weight")
	ErrContainerIsFull                        = errors.New("container is full")
	ErrNotContainer                           = errors.New("item is not a container")
	ErrEmptyLabel                             = errors.New("empty label")
	ErrNounCannotBeNil                        = errors.New("noun cannot be nil")
	ErrAdjectiveCannotBeNil                   = errors.New("adjective cannot be nil")
	ErrNounCannotBeUnderscore                 = errors.New("noun cannot be underscore")
	ErrWeightShouldBeLessOrEqualThanMaxWeight = errors.New("weight should be less or equal than max weight")
	ErrWeightCannotBeNegative                 = errors.New("weight nor max weight cannot be negative")
	ErrItemValidationFailed                   = errors.New("item validation failed")
	ErrItemCannotBeWornAndHaveLocation        = errors.New("item cannot be worn and have location")
	ErrDuplicatedItemLabel                    = errors.New("duplicated item label")
	ErrDuplicatedNounAdj                      = errors.New("duplicated noun and adjective")
)
