package database

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrNilStoreable          = errors.New("nil storeable")
	ErrDuplicatedRecord      = errors.New("duplicated record")
	ErrRecordNotFound        = errors.New("record not found")
	ErrDstMustBePointer      = errors.New("destination must be a pointer")
	ErrDstMustBePointerSlice = errors.New("destination must be a pointer to a slice")
	ErrLimitReached          = fmt.Errorf("limit reached (%d), cannot add more labels", math.MaxUint32)
	ErrNotFound              = errors.New("label not found")
	ErrOnlyLabelWithoutID    = errors.New("only labels without ID can be created")
)
