package player

import "strings"

type Store []*Player

func NewStore() Store {
	return Store{}
}

func (s Store) Get(label string) *Player {
	label = strings.ToLower(label)

	for _, p := range s {
		if p.label == label {
			return p
		}
	}

	return nil
}

func (s Store) Exists(label string) bool {
	label = strings.ToLower(label)

	for _, p := range s {
		if p.label == label {
			return true
		}
	}

	return false
}

func (s Store) GetHuman() *Player {
	for _, p := range s {
		if p.Human {
			return p
		}
	}

	return nil
}

// Set a new npc
func (s *Store) Set(n *Player) error {
	n.label = strings.ToLower(n.label)

	if s.Exists(n.label) {
		return ErrDuplicatedPlayerLabel
	}

	*s = append(*s, n)

	return nil
}
