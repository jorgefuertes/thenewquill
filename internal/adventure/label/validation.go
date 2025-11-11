package label

import (
	"errors"
	"regexp"

	"github.com/jorgefuertes/thenewquill/internal/adventure/id"
)

func (l Label) Validate(allowNoID bool) error {
	if l.ID == Undefined.ID && !allowNoID {
		return id.ErrUndefined
	}

	if !regexp.MustCompile(`^[\d\p{L}\-_\.]{1,25}$`).MatchString(l.Name) {
		return errors.Join(ErrInvalidLabelName, errors.New(l.Name))
	}

	return nil
}
