package words

import (
	"slices"
	"strings"

	"thenewquill/internal/compiler/section"
)

type Store []*Word

func NewStore() Store {
	v := Store{}

	for _, t := range wordTypes() {
		v = append(v, &Word{Label: "_", Type: t, Synonyms: make([]string, 0)})
	}

	return v
}

// Set a new word
// overwrites any existing word with the same label and type
func (s *Store) Set(label string, t WordType, synonyms ...string) *Word {
	label = strings.ToLower(label)

	if synonyms == nil {
		synonyms = make([]string, 0)
	}

	if s.Exists(t, label) {
		w := s.Get(t, label)
		w.Synonyms = synonyms

		if w.Type == Unknown {
			w.Type = t
		}

		return w
	}

	w := &Word{
		Label:    label,
		Type:     t,
		Synonyms: synonyms,
	}

	*s = append(*s, w)

	return w
}

func (s Store) Get(t WordType, labelOrSynonym string) *Word {
	labelOrSynonym = strings.ToLower(labelOrSynonym)

	for _, w := range s {
		if w.Is(labelOrSynonym) && w.Type == t {
			return w
		}
	}

	return nil
}

func (s Store) Exists(t WordType, labelOrSynonym string) bool {
	labelOrSynonym = strings.ToLower(labelOrSynonym)

	for _, w := range s {
		if w.Type != t {
			continue
		}

		if w.Label == labelOrSynonym || slices.Contains(w.Synonyms, labelOrSynonym) {
			return true
		}
	}

	return false
}

// First returns the first word of the vocabulary that matches the given label or synonym.
func (s Store) First(labelOrSynonym string) *Word {
	labelOrSynonym = strings.ToLower(labelOrSynonym)

	for _, w := range s {
		if w.Is(labelOrSynonym) {
			return w
		}
	}

	return nil
}

// FirstWithTypes returns the first word found by given types preferences
func (s Store) FirstWithTypes(labelOrSynonym string, types ...WordType) *Word {
	labelOrSynonym = strings.ToLower(labelOrSynonym)

	for _, t := range types {
		w := s.Get(t, labelOrSynonym)
		if w != nil {
			return w
		}
	}

	return nil
}

func (s Store) Len() int {
	return len(s)
}

func (s Store) Export() (section.Section, [][]string) {
	data := make([][]string, 0)

	for _, w := range s {
		data = append(data, w.export())
	}

	return section.Vars, data
}
