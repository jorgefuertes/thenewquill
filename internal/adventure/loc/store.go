package loc

import (
	"strings"
	"sync"

	"thenewquill/internal/adventure/vars"
)

type Store struct {
	mut  *sync.Mutex
	data []*Location
}

func NewStore() Store {
	return Store{mut: &sync.Mutex{}, data: make([]*Location, 0)}
}

func (s *Store) GetAll() []*Location {
	s.mut.Lock()
	defer s.mut.Unlock()

	return s.data
}

func (s *Store) Get(label string) *Location {
	if label == "" {
		return nil
	}

	label = strings.ToLower(label)

	s.mut.Lock()
	defer s.mut.Unlock()

	for _, l := range s.data {
		if l.Label == label {
			return l
		}
	}

	return nil
}

// New creates an empty location with the given label
func (s *Store) New(label string) (*Location, error) {
	l := &Location{
		Label:       label,
		Title:       Undefined,
		Description: Undefined,
		Conns:       make([]Connection, 0),
		Vars:        vars.NewStore(),
	}
	if err := s.Set(l); err != nil {
		return nil, err
	}

	return l, nil
}

// Set a new location
// overwrites any existing location with the same label
func (s *Store) Set(l *Location) error {
	if l.Label == "" {
		return ErrEmptyLabel
	}

	if l.Conns == nil {
		l.Conns = make([]Connection, 0)
	}

	if l.Vars.Regs == nil {
		l.Vars = vars.NewStore()
	}

	l.toLower()

	if s.Exists(l.Label) {
		i := s.getIndex(l.Label)
		s.mut.Lock()
		defer s.mut.Unlock()
		s.data[i] = l
		return nil
	}

	s.mut.Lock()
	defer s.mut.Unlock()
	s.data = append(s.data, l)

	return nil
}

func (s *Store) Exists(label string) bool {
	return s.getIndex(label) > -1
}

func (s *Store) Len() int {
	s.mut.Lock()
	defer s.mut.Unlock()

	return len(s.data)
}

func (s *Store) getIndex(label string) int {
	if label == "" {
		return -1
	}

	label = strings.ToLower(label)

	s.mut.Lock()
	defer s.mut.Unlock()

	for i, l := range s.data {
		if l.Label == label {
			return i
		}
	}

	return -1
}
