package character

import (
	"strings"
	"sync"
)

type Store struct {
	lock  *sync.Mutex
	chars []*Character
}

func NewStore() Store {
	return Store{lock: &sync.Mutex{}, chars: make([]*Character, 0)}
}

func (s Store) Validate() error {
	if s.GetHuman() == nil {
		return ErrNoHuman
	}

	s.lock.Lock()
	defer s.lock.Unlock()

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

func (s Store) Len() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return len(s.chars)
}

func (s Store) Get(label string) *Character {
	s.lock.Lock()
	defer s.lock.Unlock()

	label = strings.ToLower(label)

	for _, p := range s.chars {
		if p.Label == label {
			return p
		}
	}

	return nil
}

func (s Store) Exists(label string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	label = strings.ToLower(label)

	for _, p := range s.chars {
		if p.Label == label {
			return true
		}
	}

	return false
}

func (s Store) GetHuman() *Character {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, p := range s.chars {
		if p.Human {
			return p
		}
	}

	return nil
}

func (s *Store) getIndex(label string) int {
	s.lock.Lock()
	defer s.lock.Unlock()

	label = strings.ToLower(label)

	for i, p := range s.chars {
		if p.Label == label {
			return i
		}
	}

	return -1
}

// Set a new npc, if it already exists, it will be replaced
func (s *Store) Set(c *Character) error {
	c.Label = strings.ToLower(c.Label)

	if s.Exists(c.Label) {
		i := s.getIndex(c.Label)

		s.lock.Lock()
		defer s.lock.Unlock()

		s.chars[i] = c

		return nil
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.chars = append(s.chars, c)

	return nil
}
