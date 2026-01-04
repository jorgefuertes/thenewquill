package config

import (
	"fmt"
	"slices"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
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

func (s *Service) ValidateAll() []error {
	validationErrors := []error{}

	// all the required params must be set
	for _, l := range requiredParamLabels {
		if !s.Get().WithLabel(l).Exists() {
			validationErrors = append(validationErrors, fmt.Errorf("%w: param %q not found", ErrMissingConfigField, l))
		}
	}

	// check if the language is valid
	if !slices.Contains(allowedLanguages, s.GetValueOrBlank(LanguageParamLabel)) {
		validationErrors = append(
			validationErrors,
			fmt.Errorf("%w: %s=%q", ErrUnrecognizedLanguage, LanguageParamLabel, s.GetValueOrBlank(LanguageParamLabel)),
		)
	}

	// validate every param
	res := s.db.Query(database.FilterByKind(kind.Param))
	defer res.Close()

	p := &Param{}
	for res.Next(p) {
		if err := validator.Validate(p); err != nil {
			validationErrors = append(
				validationErrors,
				fmt.Errorf("%w: param #%d %s=%s", err, p.ID, s.db.GetLabelOrBlank(p.LabelID), p.V),
			)

			continue
		}

		if !IsValidLabel(s.db.GetLabelOrBlank(p.LabelID)) {
			validationErrors = append(
				validationErrors,
				fmt.Errorf("param label #%d:%s not valid", p.LabelID, s.db.GetLabelOrBlank(p.LabelID)),
			)
		}
	}

	return validationErrors
}
