package variable

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/adventure/id"
)

func (v Variable) Validate(allowNoID bool) error {
	if err := v.ID.Validate(false); err != nil && !allowNoID {
		return err
	}

	if v.ID < id.Min && !allowNoID {
		return db.ErrInvalidLabelID
	}

	return nil
}

func (s Service) ValidateAll() error {
	for _, v := range s.All() {
		if err := v.Validate(false); err != nil {
			return err
		}
	}

	return nil
}
