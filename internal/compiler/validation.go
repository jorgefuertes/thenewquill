package compiler

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func validateSection(a *adventure.Adventure, s *status.Status, k kind.Kind) error {
	validators := map[kind.Kind]func() []error{
		kind.Param:     a.Config.ValidateAll,
		kind.Word:      a.Words.ValidateAll,
		kind.Message:   a.Messages.ValidateAll,
		kind.Variable:  a.Variables.ValidateAll,
		kind.Item:      a.Items.ValidateAll,
		kind.Character: a.Characters.ValidateAll,
		kind.Location:  a.Locations.ValidateAll,
	}

	if s.HasRunValidator(k) {
		return nil
	}

	s.FlagValidator(k)

	validationErrors := validators[k]()
	if len(validationErrors) > 0 {
		e := cerr.ErrValidation.WithSection(k).
			WithStack(s.Stack).
			WithFilename(s.CurrentFilename())

		for _, e2 := range validationErrors {
			e = e.AddErr(e2)
		}

		return e
	}

	return nil
}
