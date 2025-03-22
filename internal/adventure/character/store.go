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

func (s *Store) Validate() error {
	if s.GetHuman() == nil {
		return ErrNoHuman
	}

	s.mut.Lock()
	defer s.mut.Unlock()

	humans := 0
	for _, p := range s.chars {
		if p.Human {
			humans++
		}
	}

	if humans > 1 {
		return ErrOnlyOneHuman
	}

	return nil
}

func (s *Store) Len() int {
	s.mut.Lock()
	defer s.mut.Unlock()

	return len(s.chars)
}

func (s *Store) CreateEmpty(label string) *Character {
	c := New(label, nil, nil)
	s.Set(c)

	return c
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
	s.mut.Lock()
	defer s.mut.Unlock()

	label = strings.ToLower(label)

	for _, p := range s.chars {
		if p.Label == label {
			return true
		}
	}

	return false
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

// Set a new npc, if it already exists, it will be replaced
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

	return
}
