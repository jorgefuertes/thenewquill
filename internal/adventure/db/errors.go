package db

import (
	"errors"
	"fmt"
)

var (
	ErrDuplicatedRecord          = errors.New("duplicated record")
	ErrRecordNotFound            = errors.New("record not found")
	ErrDstMustBePointer          = errors.New("destination must be a pointer")
	ErrDstMustBePointerSlice     = errors.New("destination must be a pointer to a slice")
	ErrLimitReached              = fmt.Errorf("limit reached (%d), cannot add more labels", Limit)
	ErrNotFound                  = errors.New("label not found")
	ErrCannotCastFromStoreable   = errors.New("cannot cast from storeable")
	ErrUndefinedLabel            = errors.New("undefined label")
	ErrInvalidLabelName          = errors.New("invalid label name")
	ErrInvalidLabelID            = errors.New("invalid label ID")
	ErrKindCannotBeNone          = errors.New("kind cannot be none")
	ErrSubKindMustBeDefined      = errors.New("subkind must be defined")
	ErrCannotCreateWithDefinedID = errors.New("cannot create with defined ID")
)
