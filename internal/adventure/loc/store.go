package loc

import (
	"strings"

	"thenewquill/internal/compiler/section"
)

type Store []*Location

func NewStore() Store {
	return Store{}
}

func (s Store) Get(label string) *Location {
	label = strings.ToLower(label)

	for _, l := range s {
		if l.Label == label {
			return l
		}
	}

	return nil
}

// Set a new location
// overwrites any existing location with the same label
func (s *Store) Set(label, title, desc string) *Location {
	label = strings.ToLower(label)

	if existent := s.Get(label); existent != nil {
		existent.Title = title
		existent.Description = desc

		return existent
	}

	l := New(label, title, desc)
	*s = append(*s, l)

	return l
}

func (s Store) Exists(label string) bool {
	label = strings.ToLower(label)

	for _, l := range s {
		if l.Label == label {
			return true
		}
	}

	return false
}

func (s Store) Len() int {
	return len(s)
}

func (s Store) Export() (section.Section, [][]string) {
	data := make([][]string, 0)

	for _, l := range s {
		data = append(data, l.export())
	}

	return section.Locs, data
}
