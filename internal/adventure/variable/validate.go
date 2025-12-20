package variable

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/internal/database/primitive"
)

func (v Variable) Validate(allowNoID bool) error {
	if err := v.ID.ValidateID(false); err != nil && !allowNoID {
		return err
	}

	if v.ID < primitive.MinID && !allowNoID {
		return primitive.ErrInvalidID
	}

	return nil
}

func (s Service) ValidateAll() error {
	res := s.db.Query(database.FilterByKind(kind.Variable))
	defer res.Close()

	var v Variable
	for res.Next(&v) {
		if err := v.Validate(false); err != nil {
			return err
		}
	}

	return nil
}
