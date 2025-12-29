package config

import (
	"errors"
	"fmt"
	"slices"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/pkg/log"
	"github.com/jorgefuertes/thenewquill/pkg/validator"
)

const (
	TitleParamLabel       = "title"
	AuthorParamLabel      = "author"
	DescriptionParamLabel = "description"
	VersionParamLabel     = "version"
	DateParamLabel        = "date"
	LanguageParamLabel    = "language"
)

var (
	allowedLanguages   = []string{"en", "es"}
	allowedParamLabels = []string{
		TitleParamLabel,
		AuthorParamLabel,
		DescriptionParamLabel,
		VersionParamLabel,
		DateParamLabel,
		LanguageParamLabel,
	}
	requiredParamLabels = []string{TitleParamLabel, AuthorParamLabel, VersionParamLabel, LanguageParamLabel}
)

func IsValidLabel(label string) bool {
	return slices.Contains(allowedParamLabels, label)
}

func GetAllowedParamLabels() []string {
	return allowedParamLabels
}

func (s *Service) ValidateAll() error {
	seen := []string{}

	res := s.db.Query(database.FilterByKind(kind.Param))
	defer res.Close()

	p := &Param{}
	for res.Next(p) {
		log.Debug("[CONFIG:VALIDATION:ALL] Validating %q->%v", s.db.GetLabelOrBlank(p.LabelID), p)

		if err := validator.Validate(p); err != nil {
			return err
		}

		label, err := s.db.GetLabel(p.LabelID)
		if err != nil {
			return err
		}

		if !IsValidLabel(label) {
			return ErrUnrecognizedConfigField
		}

		if label == "language" && !slices.Contains(allowedLanguages, fmt.Sprintf("%v", p.V)) {
			return ErrUnrecognizedLanguage
		}

		seen = append(seen, label)
		log.Debug("[CONFIG:VALIDATION:ALL] Found config field %q", label)
	}

	// check required
	for _, l := range requiredParamLabels {
		if !slices.Contains(seen, l) {
			return errors.Join(ErrMissingConfigField, fmt.Errorf("label %q not found", l))
		}
	}

	return nil
}
