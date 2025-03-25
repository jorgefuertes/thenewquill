package character

import (
	"strings"
	"sync"
)

type Store struct {
	mut   *sync.Mutex
	chars []*Character
}

func NewStore() Store {
	return Store{mut: &sync.Mutex{}, chars: make([]*Character, 0)}
}

func (s *Store) Len() int {
	s.mut.Lock()
	defer s.mut.Unlock()

	return len(s.chars)
}

// Set a npc, if it already exists, it will be replaced
func (s *Store) Set(c *Character) {
	c.Label = strings.ToLower(c.Label)

	if s.Exists(c.Label) {
		i := s.getIndex(c.Label)

		s.mut.Lock()
		defer s.mut.Unlock()

		s.chars[i] = c

		return
	}

	s.mut.Lock()
	defer s.mut.Unlock()

	s.chars = append(s.chars, c)
}

// New creates an empty character with the given label
func (s *Store) New(label string) *Character {
	c := &Character{Label: label}
	s.Set(c)

	return c
}

func (s *Store) GetAll() []*Character {
	s.mut.Lock()
	defer s.mut.Unlock()

	return s.chars
}

func (s *Store) Get(label string) *Character {
	if label == "" {
		return nil
	}

	s.mut.Lock()
	defer s.mut.Unlock()

	label = strings.ToLower(label)

	for _, p := range s.chars {
		if p.Label == label {
			return p
		}
	}

	return nil
}

func (s *Store) Exists(label string) bool {
	return s.getIndex(label) > -1
}

func (s *Store) GetHuman() *Character {
	s.mut.Lock()
	defer s.mut.Unlock()

	for _, p := range s.chars {
		if p.Human {
			return p
		}
	}

	return nil
}

func (s *Store) getIndex(label string) int {
	s.mut.Lock()
	defer s.mut.Unlock()

	label = strings.ToLower(label)

	for i, p := range s.chars {
		if p.Label == label {
			return i
		}
	}

	return -1
}
