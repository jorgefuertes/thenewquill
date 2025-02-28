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
	if s.Exists(m.Type, m.Label) && !m.IsPluralized() {
		return ErrMsgAlreadyExists
	}

	if m.IsPluralized() && s.Exists(m.Type, m.Label) {
		old := s.Get(m.Type, m.Label)
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

func (s Store) Exists(t MsgType, label string) bool {
	for _, msg := range s {
		if msg.Type == t && msg.Label == label {
			return true
		}
	}

	return false
}

func (s Store) Get(t MsgType, label string) *Msg {
	for _, m := range s {
		if m.Type == t && m.Label == label {
			return m
		}
	}

	return nil
}

func (s Store) Len() int {
	return len(s)
}

func (s Store) Check() error {
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
