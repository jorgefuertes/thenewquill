package config

import (
	"errors"
	"fmt"
	"slices"

	"github.com/jorgefuertes/thenewquill/internal/adventure/database"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

var allowedLanguages = []string{"en", "es"}

func (v Param) Validate(allowNoID bool) error {
	if err := v.ID.ValidateID(false); err != nil && !allowNoID {
		return err
	}

	if fmt.Sprintf("%v", v.V) == "" {
		return ErrValueIsEmpty
	}

	return nil
}

func ValidateAll(db *database.DB) error {
	seen := []primitive.Label{}

	res := db.Query(database.FilterByKind(kind.Param))

	p := Param{}
	for res.Next(p) {
		if err := p.Validate(false); err != nil {
			return err
		}

		l, err := db.GetLabel(p.LabelID)
		if err != nil {
			return err
		}

		if !isLabelAllowed(l) {
			return ErrUnrecognizedConfigField
		}

		if l == "lang" && !slices.Contains(allowedLanguages, fmt.Sprintf("%v", p.V)) {
			return ErrUnrecognizedLanguage
		}

		seen = append(seen, l)
	}

	// check required
	for _, allowed := range allowedFieldLabels {
		if !slices.Contains(seen, allowed.label) && allowed.required {
			return errors.Join(ErrMissingConfigField, fmt.Errorf("label %q not found", allowed.label))
		}
	}

	return nil
}
