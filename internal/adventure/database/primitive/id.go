package primitive

import (
	"fmt"
	"strconv"
)

const (
	UndefinedID ID = 0
	MinID       ID = 1
)

type ID uint32

func IDFromString(s string) ID {
	i, _ := strconv.ParseUint(s, 10, 32)

	return ID(i)
}

func (id ID) ToString() string {
	return fmt.Sprintf("%d", id)
}

// Validate checks if the ID is valid, allowSpecial allows _ and * to be valid IDs
func (id ID) ValidateID(allowUnderMin bool) error {
	if id.IsUndefinedID() {
		return ErrUndefinedID
	}

	if id < MinID && !allowUnderMin {
		return ErrInvalidID
	}

	return nil
}

func (id ID) IsDefinedID() bool {
	return id != UndefinedID
}

func (id ID) IsUndefinedID() bool {
	return id == UndefinedID
}
