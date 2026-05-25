package character

import (
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/pkg/validator"
)

func (s *Service) ValidateAll() []error {
	validationErrors := []error{}

	humansCount := s.db.Query(database.FilterByKind(kind.Character), database.NewFilter("Human", database.Equal, true)).
		Count()
	if humansCount == 0 {
		validationErrors = append(validationErrors, ErrNoHuman)
	} else if humansCount > 1 {
		validationErrors = append(validationErrors, ErrOnlyOneHuman)
	}

	chars := s.db.Query(database.FilterByKind(kind.Character))
	defer chars.Close()

	c := &Character{}
	for chars.Next(c) {
		if err := validator.Validate(c); err != nil {
			validationErrors = append(
				validationErrors,
				fmt.Errorf("%w: character %d %q", err, c.ID, s.db.GetLabelOrBlank(c.LabelID)),
			)
		}

		// same label as item its not allowed
		q := s.db.Query(database.FilterByKind(kind.Item), database.FilterByLabelID(c.LabelID))
		if q.Exists() {
			validationErrors = append(
				validationErrors,
				fmt.Errorf("%w: character %d %q", ErrItemHasSameLabel, c.ID, s.db.GetLabelOrBlank(c.LabelID)),
			)
		}

		// same noun as item its not allowed
		q = s.db.Query(database.FilterByKind(kind.Item), database.NewFilter("NounID", database.Equal, c.NounID))
		if q.Exists() {
			validationErrors = append(
				validationErrors,
				fmt.Errorf("%w: character %d %q", ErrItemHasSameNoun, c.ID, s.db.GetLabelOrBlank(c.LabelID)),
			)
		}
	}

	return validationErrors
}
