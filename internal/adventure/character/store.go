package character

import "strings"

type Store []*Character

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
