package variable

import (
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/pkg/validator"
)

func (v Variable) Validate() error {
	if err := validator.Validate(v); err != nil {
		return err
	}

	return nil
}

func (s Service) ValidateAll() []error {
	validationErrors := []error{}

	res := s.db.Query(database.FilterByKind(kind.Variable))
	defer res.Close()

	var v Variable
	for res.Next(&v) {
		if err := v.Validate(); err != nil {
			validationErrors = append(
				validationErrors,
				fmt.Errorf("%w: variable %d %q=%s", err, v.ID, s.db.GetLabelOrBlank(v.LabelID), v.Value),
			)
		}
	}

	return validationErrors
}
