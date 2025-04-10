package msg

import (
	"errors"
	"fmt"
)

type Store []*Msg

func NewStore() Store {
	return Store{}
}

func (s *Store) Set(m *Msg) error {
	if s.Exists(m.Label) && !m.IsPluralized() {
		return ErrMsgAlreadyExists
	}

	if m.IsPluralized() && s.Exists(m.Label) {
		old := s.Get(m.Label)
		if !old.IsPluralized() {
			return ErrMsgAlreadyExists
		}

		for i, text := range m.Plurals {
			if text != "" {
				if old.Plurals[i] == "" {
					old.Plurals[i] = text

					return nil
				}

				return ErrMsgAlreadyExists
			}
		}

		return nil
	}

	*s = append(*s, m)

	return nil
}

func (s Store) Exists(label string) bool {
	for _, msg := range s {
		if msg.Label == label {
			return true
		}
	}

	return false
}

func (s Store) Get(label string) *Msg {
	for _, m := range s {
		if m.Label == label {
			return m
		}
	}

	return nil
}

func (s Store) Len() int {
	return len(s)
}

func (s Store) Validate() error {
	for _, m := range s {
		if m.IsPluralized() {
			for i, text := range m.Plurals {
				if text == "" {
					return errors.Join(ErrMsgPluralEmpty, fmt.Errorf("label: %s missing: %s", m.Label, PluralNames[i]))
				}
			}

			continue
		}

		if m.Text == "" {
			return errors.Join(ErrMsgEmpty, fmt.Errorf("label: %s", m.Label))
		}
	}

	return nil
}
