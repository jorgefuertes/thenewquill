package status

import "github.com/jorgefuertes/thenewquill/internal/adventure/kind"

func (s *Status) FlagValidator(k kind.Kind) {
	if !s.HasRunValidator(k) {
		s.runnedValidators = append(s.runnedValidators, k)
	}
}

func (s *Status) FlagReplacer(k kind.Kind) {
	if !s.HasRunReplacer(k) {
		s.runnedReplacers = append(s.runnedReplacers, k)
	}
}

func (s *Status) HasRunValidator(k kind.Kind) bool {
	for _, v := range s.runnedValidators {
		if v == k {
			return true
		}
	}

	return false
}

func (s *Status) HasRunReplacer(k kind.Kind) bool {
	for _, v := range s.runnedReplacers {
		if v == k {
			return true
		}
	}

	return false
}
