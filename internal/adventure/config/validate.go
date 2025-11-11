package config

import (
	"errors"
	"fmt"
	"slices"
)

var allowedLanguages = []string{"en", "es"}

func (v Param) Validate(allowNoID bool) error {
	if err := v.ID.Validate(false); err != nil && !allowNoID {
		return err
	}

	if fmt.Sprintf("%v", v.V) == "" {
		return ErrValueIsEmpty
	}

	return nil
}

func (s *Service) ValidateAll() error {
	seen := []string{}

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

		seen = append(seen, label.Name)
	}

	// check required
	for _, allowed := range allowedFields {
		if !slices.Contains(seen, allowed.labelName) && allowed.required {
			return errors.Join(ErrMissingConfigField, fmt.Errorf("label %q not found", allowed.labelName))
		}
	}

	return nil
}
