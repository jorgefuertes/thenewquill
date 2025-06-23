package item

import (
	"errors"
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

func (i Item) Validate(allowNoID db.Allow) error {
	if err := i.ID.Validate(db.DontAllowSpecial); err != nil && !allowNoID {
		return fmt.Errorf("ID %q: %w", i.ID, err)
	}

	if err := i.NounID.Validate(db.DontAllowSpecial); err != nil {
		return fmt.Errorf("noun ID %q: %w", i.NounID, err)
	}

	if err := i.AdjectiveID.Validate(db.AllowSpecial); err != nil {
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
	for _, i := range s.All() {
		label, err := s.db.GetLabel(i.ID)
		if err != nil {
			return errors.Join(err, fmt.Errorf("item ID: %d", i.ID))
		}

		if err := i.Validate(db.DontAllowNoID); err != nil {
			return errors.Join(
				ErrItemValidationFailed,
				err,
				fmt.Errorf("label %d: %s", label.ID, label.Name),
			)
		}

		for _, i2 := range s.All() {
			if i.ID == i2.ID {
				continue
			}

			if i.NounID == i2.NounID && i.AdjectiveID == i2.AdjectiveID {
				label2, err := s.db.GetLabel(i2.ID)
				if err != nil {
					return errors.Join(
						err,
						fmt.Errorf("items %q and %d, noun %d and adjective %d", label, i2.ID, i.NounID, i.AdjectiveID),
					)
				}
				return errors.Join(
					ErrDuplicatedNounAdj,
					fmt.Errorf("items %q and %q: noun %d and adjective %d", label, label2, i.NounID, i.AdjectiveID),
				)
			}
		}
	}

	return nil
}
