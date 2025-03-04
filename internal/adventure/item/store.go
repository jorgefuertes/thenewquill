package item

import (
	"strings"

	"thenewquill/internal/adventure/words"
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
		if item.Label == label {
			return item
		}
	}

	return nil
}

func (s Store) Exists(label string) bool {
	for _, item := range s {
		if item.Label == label {
			return true
		}
	}

	return false
}

func (s Store) ExistsNounAdj(noun, adjective *words.Word) bool {
	for _, item := range s {
		if item.Noun.IsEqual(noun) && item.Adjective.IsEqual(adjective) {
			return true
		}
	}

	return false
}

func (s *Store) Set(newItem *Item) error {
	if s.Exists(newItem.Label) {
		return ErrDuplicateLabel
	}

	if newItem.Label == "" {
		return ErrEmptyLabel
	}

	newItem.Label = strings.ToLower(newItem.Label)

	if newItem.Noun == nil {
		return ErrNounCannotBeNil
	}

	if newItem.Noun.Label == "_" {
		return ErrNounCannotBeUnderscore
	}

	*s = append(*s, newItem)

	return nil
}
