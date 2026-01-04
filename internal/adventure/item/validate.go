package item

import (
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

func (s *Service) ValidateAll() []error {
	validationErrors := []error{}

	res := s.db.Query(database.FilterByKind(kind.Item))
	defer res.Close()

	i := Item{}
	for res.Next(&i) {
		if err := i.Validate(); err != nil {
			validationErrors = append(
				validationErrors,
				fmt.Errorf("%w: item #%d:%q", err, i.ID, s.db.GetLabelOrBlank(i.LabelID)),
			)
		}

		if s.TotalWeight(i) > i.MaxWeight {
			validationErrors = append(
				validationErrors,
				fmt.Errorf(
					"%w: item #%d:%q, weight %d/%d",
					ErrContainerCantCarrySoMuch,
					i.ID,
					s.db.GetLabelOrBlank(i.LabelID),
					s.TotalWeight(i),
					i.MaxWeight,
				),
			)
		}

		if s.db.Query(
			database.NewFilter("ID", database.NotEqual, i.ID),
			database.FilterByKind(kind.Item),
			database.NewFilter("NounID", database.Equal, i.NounID),
			database.NewFilter("AdjectiveID", database.Equal, i.AdjectiveID),
		).Exists() {
			validationErrors = append(
				validationErrors,
				fmt.Errorf(
					"%w: item #%d:%q noun: #%d:%q adj: #%d:%q",
					ErrDuplicatedNounAdj,
					i.ID,
					s.db.GetLabelOrBlank(i.LabelID),
					i.NounID,
					s.db.GetLabelFromRecordOrBlank(i.NounID),
					i.AdjectiveID,
					s.db.GetLabelFromRecordOrBlank(i.AdjectiveID),
				),
			)
		}
	}

	return validationErrors
}
