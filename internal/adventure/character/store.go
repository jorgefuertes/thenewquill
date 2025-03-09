package character

import (
	"strings"

	"thenewquill/internal/compiler/section"
)

type Store []*Character

func (s Store) Validate() error {
	if s.GetHuman() == nil {
		return ErrNoHuman
	}

	humans := 0
	for _, p := range s {
		if p.Human {
			humans++
		}
	}

	if humans > 1 {
		return ErrOnlyOneHuman
	}

	return nil
}

func NewStore() Store {
	return Store{}
}

func (s Store) Len() int {
	return len(s)
}

func (s Store) Get(label string) *Character {
	label = strings.ToLower(label)

	for _, p := range s {
		if p.Label == label {
			return p
		}
	}

	return nil
}

func (s Store) Exists(label string) bool {
	label = strings.ToLower(label)

	for _, p := range s {
		if p.Label == label {
			return true
		}
	}

	return false
}

func (s Store) GetHuman() *Character {
	for _, p := range s {
		if p.Human {
			return p
		}
	}

	return nil
}

// Set a new npc
func (s *Store) Set(n *Character) error {
	n.Label = strings.ToLower(n.Label)

	if s.Exists(n.Label) {
		return ErrDuplicatedPlayerLabel
	}

	*s = append(*s, n)

	return nil
}

func (s Store) Export() (section.Section, [][]string) {
	data := make([][]string, 0)

	for _, c := range s {
		data = append(data, c.export())
	}

	return section.Chars, data
}
