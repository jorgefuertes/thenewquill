package database

import "errors"

var (
	ErrRecordNotFound             = errors.New("record not found")
	ErrLabelNotFound              = errors.New("label not found")
	ErrEntityIsNotPointer         = errors.New("entity should be a pointer")
	ErrEntityIsNotPointerToStruct = errors.New("entity should be a pointer to struct")
	ErrMissingIDField             = errors.New("missing 'ID' field")
	ErrIDFieldIsNotUint32         = errors.New("ID field is not uint32")
	ErrMissingLabelIDField        = errors.New("missing 'LabelID' field")
	ErrLabelFieldIsNotUint32      = errors.New("LabelID field is not uint32")
	ErrCannotInferKind            = errors.New("cannot infer entity kind")
	ErrDatabaseIsFull             = errors.New("database is full")
	ErrLabelsAreFull              = errors.New("labels are full")
	ErrCannotSetIDField           = errors.New("cannot set 'ID' field")
	ErrIDFieldIsNotZero           = errors.New("ID field is not zero")
	ErrInvalidLabel               = errors.New("invalid label")
	ErrDatabaseIsFrozen           = errors.New("database is frozen, can't create new records")
	ErrMissingIDToUpdate          = errors.New("missing ID to update")
	ErrWrongUpdateKind            = errors.New("cannot update record with different kind entity")
)
