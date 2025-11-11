package id

import (
	"fmt"
	"strconv"
)

const (
	Undefined ID = 0
	Min       ID = 4
)

type ID uint32

func FromString(s string) ID {
	i, _ := strconv.ParseUint(s, 10, 32)

	return ID(i)
}

func (i ID) ToString() string {
	return fmt.Sprintf("%d", i)
}

// Validate checks if the ID is valid, allowSpecial allows _ and * to be valid IDs
func (i ID) Validate(allowUnderMin bool) error {
	if i == Undefined {
		return ErrUndefined
	}

	if i < Min && !allowUnderMin {
		return ErrInvalid
	}

	return nil
}

func (i ID) IsDefined() bool {
	return i != Undefined
}
