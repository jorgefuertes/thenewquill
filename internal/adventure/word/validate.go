package word

import (
	"fmt"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/internal/lang"
	"github.com/jorgefuertes/thenewquill/pkg/validator"
)

func (w Word) Validate() error {
	return validator.Validate(w)
}

func (s *Service) ValidateAll() []error {
	validationErrors := []error{}

	words := s.db.Query(database.FilterByKind(kind.Word))
	defer words.Close()

	var w Word
	for words.Next(&w) {
		label := s.db.GetLabelOrBlank(w.LabelID)

		// ignore underscore and asterisk
		if label == database.LabelUnderscore || label == database.LabelAsterisk {
			continue
		}

		if err := w.Validate(); err != nil {
			validationErrors = append(
				validationErrors,
				err,
				fmt.Errorf("ID #%d LABEL #%d:%s", w.ID, w.LabelID, s.db.GetLabelOrBlank(w.LabelID)),
				fmt.Errorf("synonyms: %s", strings.Join(w.Synonyms, ", ")),
			)
		}

		// duplicated label and type
		if s.Get().WithNoID(w.ID).WithType(w.Type).WithLabelID(w.LabelID).Exists() {
			validationErrors = append(
				validationErrors,
				fmt.Errorf("%w: duplicated label %q and type %q", ErrDuplicatedLabel, label, w.Type),
			)
		}

		// duplicated synonym in the same type
		for _, syn := range w.Synonyms {
			if s.Get().WithNoID(w.ID).WithType(w.Type).WithSynonym(syn).Exists() {
				w2, _ := s.Get().WithNoID(w.ID).WithType(w.Type).WithSynonym(syn).First()

				validationErrors = append(
					validationErrors,
					fmt.Errorf(
						"%w: %q between %s %q and %s %q",
						ErrDuplicatedSyn,
						syn,
						w.Type,
						label,
						w2.Type,
						s.db.GetLabelOrBlank(w2.LabelID),
					),
				)
			}
		}
	}

	// check for required verbs
	l := s.GetLang()

	if s.GetDefaultVerbSyns(lang.Lang(l), lang.Go) == nil {
		validationErrors = append(validationErrors, ErrMissingGoVerb)
	}

	if s.GetDefaultVerbSyns(lang.Lang(l), lang.Talk) == nil {
		validationErrors = append(validationErrors, ErrMissingTalkVerb)
	}

	if s.GetDefaultVerbSyns(lang.Lang(l), lang.Examine) == nil {
		validationErrors = append(validationErrors, ErrMissingExamineVerb)
	}

	return validationErrors
}
