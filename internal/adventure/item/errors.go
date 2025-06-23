package item

import "errors"

var (
	ErrItemNotContained                       = errors.New("item is not inside any container")
	ErrItemAlreadyContained                   = errors.New("item already contained elsewhere")
	ErrContainerCantCarrySoMuch               = errors.New("container can't carry so much weight")
	ErrNotContainer                           = errors.New("item is not a container")
	ErrNounCannotBeNil                        = errors.New("noun cannot be nil")
	ErrAdjectiveCannotBeNil                   = errors.New("adjective cannot be nil")
	ErrNounCannotBeUnderscore                 = errors.New("noun cannot be underscore")
	ErrWeightShouldBeLessOrEqualThanMaxWeight = errors.New("weight should be less or equal than max weight")
	ErrWeightCannotBeNegative                 = errors.New("weight nor max weight cannot be negative")
	ErrItemValidationFailed                   = errors.New("item validation failed")
	ErrItemCannotBeWornAndHaveLocation        = errors.New("item cannot be worn and have location")
	ErrDuplicatedItemLabel                    = errors.New("duplicated item label")
	ErrDuplicatedNounAdj                      = errors.New("duplicated noun and adjective")
	ErrItemIsNotWearableButIsWorn             = errors.New("item is not wearable but is worn")
	ErrCannotAssertIntoItem                   = errors.New("cannot assert into item")
	ErrInvalidTo                              = errors.New("invalid to, should be a container, character or location")
)
