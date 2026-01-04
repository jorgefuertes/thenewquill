package compiler

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func postReplaceSection(a *adventure.Adventure, s *status.Status, k kind.Kind) error {
	replacers := map[kind.Kind]func() error{
		kind.Location: a.Locations.PostReplace,
	}

	if _, ok := replacers[k]; !ok {
		return nil
	}

	if s.HasRunReplacer(k) {
		return nil
	}

	s.FlagReplacer(k)

	return replacers[k]()
}
