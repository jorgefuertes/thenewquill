package config

import (
	"fmt"
	"slices"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

var allowedLanguages = []string{"en", "es"}

func (v Value) Validate(allowNoID db.Allow) error {
	if err := v.ID.Validate(db.DontAllowSpecial); err != nil && !allowNoID {
		return err
	}

	if fmt.Sprintf("%v", v.V) == "" {
		return ErrValueIsEmpty
	}

	return nil
}

func (s *Service) ValidateAll() error {
	for _, v := range s.All() {
		if err := v.Validate(false); err != nil {
			return err
		}

		label, err := s.db.GetLabel(v.ID)
		if err != nil {
			return err
		}

		if !isKeyAllowed(label.Name) {
			return ErrUnrecognizedConfigField
		}

		if label.Name == "lang" && !slices.Contains(allowedLanguages, fmt.Sprintf("%v", v.V)) {
			return ErrUnrecognizedLanguage
		}
	}

	return nil
}
