package item

import (
	"thenewquill/internal/adventure/voc"
)

type Store []*Item

func NewStore() Store {
	return Store{}
}

func (s Store) Len() int {
	return len(s)
}

func (s Store) Get(label string) *Item {
	for _, item := range s {
		if item.label == label {
			return item
		}
	}

	return nil
}

func (s Store) Exists(label string) bool {
	for _, item := range s {
		if item.label == label {
			return true
		}
	}

	return false
}

func (s Store) ExistsNounAdj(noun, adjective *voc.Word) bool {
	for _, item := range s {
		if item.noun.IsEqual(noun) && item.adjective.IsEqual(adjective) {
			return true
		}
	}

	return false
}

func (s *Store) Set(newItem *Item) error {
	if s.Exists(newItem.label) {
		return ErrDuplicateLabel
	}

	if newItem.label == "" {
		return ErrEmptyLabel
	}

	if newItem.noun == nil {
		return ErrNounCannotBeNil
	}

	if newItem.noun.Label == "_" {
		return ErrNounCannotBeUnderscore
	}

	*s = append(*s, newItem)

	return nil
}
