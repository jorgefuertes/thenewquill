package item

import (
	"errors"
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

func (i Item) Validate() error {
	if i.NounID == db.UndefinedLabel.ID {
		return ErrNounCannotBeNil
	}

	if i.NounID == db.UnderscoreLabel.ID {
		return ErrNounCannotBeUnderscore
	}

	if i.AdjectiveID == db.UndefinedLabel.ID {
		return ErrAdjectiveCannotBeNil
	}

	if i.Weight > i.MaxWeight {
		return ErrWeightShouldBeLessOrEqualThanMaxWeight
	}

	if i.Weight < 0 || i.MaxWeight < 0 {
		return ErrWeightCannotBeNegative
	}

	if i.WornBy != db.UndefinedLabel.ID && !i.IsWearable {
		return ErrItemIsNotWearableButIsWorn
	}

	if i.WornBy != db.UndefinedLabel.ID && i.LocationID != db.UndefinedLabel.ID {
		return ErrItemCannotBeWornAndHaveLocation
	}

	return nil
}

func ValidateAll(d *db.DB) error {
	items := make([]Item, 0)

	if err := d.GetByKindAs(db.Items, db.NoSubKind, &items); err != nil {
		return err
	}

	for _, obj := range items {
		objLabel, err := d.GetLabel(obj.GetID())
		if err != nil {
			return errors.Join(err, fmt.Errorf("item ID: %d", obj.GetID()))
		}

		if err := obj.Validate(); err != nil {
			return errors.Join(
				ErrItemValidationFailed,
				err,
				fmt.Errorf("label %d: %s", objLabel.ID, objLabel.Name),
			)
		}

		items2 := make([]Item, 0)
		if err := d.GetByKindAs(db.Items, db.NoSubKind, &items2); err != nil {
			return err
		}

		for _, obj2 := range items2 {
			if obj.GetID() == obj2.GetID() {
				continue
			}

			if obj.NounID == obj2.NounID && obj.AdjectiveID == obj2.AdjectiveID {
				return errors.Join(ErrDuplicatedNounAdj, errors.New(d.GetLabelName(obj.GetID())))
			}
		}
	}

	return nil
}
