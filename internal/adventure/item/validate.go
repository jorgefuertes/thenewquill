package item

import (
	"errors"
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/internal/database/primitive"
)

func (i Item) Validate(allowNoID bool) error {
	if err := i.ID.ValidateID(false); err != nil && !allowNoID {
		return fmt.Errorf("ID %q: %w", i.ID, err)
	}

	if err := i.NounID.ValidateID(false); err != nil {
		return fmt.Errorf("noun ID %q: %w", i.NounID, err)
	}

	if err := i.AdjectiveID.ValidateID(true); err != nil {
		return fmt.Errorf("adjective ID %q: %w", i.AdjectiveID, err)
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

	i := &Item{}
	for res.Next(i) {
		l, err := s.db.GetLabel(i.ID)
		if err != nil {
			return errors.Join(err, fmt.Errorf("item %q: %d", l, i.ID))
		}

		if err := i.Validate(false); err != nil {
			return errors.Join(
				ErrItemValidationFailed,
				err,
				fmt.Errorf("item %q: %d", l, i.ID),
			)
		}

		if err := s.findDuplicatedNounAdj(l, i); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) findDuplicatedNounAdj(l primitive.Label, i *Item) error {
	res2 := s.db.Query(database.FilterByKind(kind.Item))
	defer res2.Close()

	i2 := &Item{}
	for res2.Next(i2) {
		if i.ID == i2.ID {
			continue
		}

		if i.NounID == i2.NounID && i.AdjectiveID == i2.AdjectiveID {
			l2, err := s.db.GetLabel(i2.ID)
			if err != nil {
				return errors.Join(
					err,
					fmt.Errorf("items %q and %d, noun %d and adjective %d", l, i2.ID, i.NounID, i.AdjectiveID),
				)
			}

			return errors.Join(
				ErrDuplicatedNounAdj,
				fmt.Errorf("items %q and %q: noun %d and adjective %d", l, l2, i.NounID, i.AdjectiveID),
			)
		}
	}

	return nil
}
