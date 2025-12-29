package item

import (
	"errors"
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/pkg/validator"
)

func (i Item) Validate() error {
	if err := validator.Validate(i); err != nil {
		return err
	}

	if i.Weight > i.MaxWeight {
		return ErrWeightShouldBeLessOrEqualThanMaxWeight
	}

	if i.Weight < 0 || i.MaxWeight < 0 {
		return ErrWeightCannotBeNegative
	}

	return nil
}

func (s *Service) ValidateAll() error {
	res := s.db.Query(database.FilterByKind(kind.Item))
	defer res.Close()

	i := Item{}
	for res.Next(&i) {
		if err := i.Validate(); err != nil {
			return errors.Join(fmt.Errorf("item %q", s.db.GetLabelOrBlank(i.ID)), err)
		}

		if s.TotalWeight(i) > i.MaxWeight {
			return ErrContainerCantCarrySoMuch
		}

		if s.db.Query(
			database.NewFilter("ID", database.NotEqual, i.ID),
			database.FilterByKind(kind.Item),
			database.NewFilter("NounID", database.Equal, i.NounID),
			database.NewFilter("AdjectiveID", database.Equal, i.AdjectiveID),
		).Exists() {
			return ErrDuplicatedNounAdj
		}
	}

	return nil
}
