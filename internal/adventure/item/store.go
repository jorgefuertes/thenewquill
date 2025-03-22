package item

import (
	"errors"
	"strings"
	"sync"

	"thenewquill/internal/adventure/words"
)

type Store struct {
	lock  *sync.Mutex
	items []*Item
}

func NewStore() Store {
	return Store{
		lock:  &sync.Mutex{},
		items: make([]*Item, 0),
	}
}

func (s *Store) Len() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return len(s.items)
}

func (s *Store) getIndex(label string) int {
	s.lock.Lock()
	defer s.lock.Unlock()

	for i, item := range s.items {
		if item.Label == label {
			return i
		}
	}

	return -1
}

func (s *Store) Get(label string) *Item {
	if label == "" {
		return nil
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	for _, item := range s.items {
		if item.Label == label {
			return item
		}
	}

	return nil
}

func (s *Store) Exists(label string) bool {
	i := s.getIndex(label)

	return i > -1
}

func (s Store) ExistsNounAdj(noun, adjective *words.Word) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, item := range s.items {
		if item.Noun.IsEqual(noun) && item.Adjective.IsEqual(adjective) {
			return true
		}
	}

	return false
}

func (s *Store) CreateEmpty(label string) *Item {
	i := New(label, nil, nil)
	s.Set(i)

	return i
}

// Set a new item, if it already exists, it will be replaced
func (s *Store) Set(newItem *Item) error {
	if newItem.Label == "" {
		return ErrEmptyLabel
	}

	newItem.Label = strings.ToLower(newItem.Label)

	idx := s.getIndex(newItem.Label)
	if idx > -1 {
		s.lock.Lock()
		defer s.lock.Unlock()

		s.items[idx] = newItem

		return nil
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.items = append(s.items, newItem)

	return nil
}

func (s *Store) Validate() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, item := range s.items {
		if err := item.validate(); err != nil {
			return errors.Join(ErrItemValidationFailed, err, errors.New(item.Label))
		}
	}

	for i, item := range s.items {
		for i2, item2 := range s.items {
			if i == i2 {
				continue
			}

			if i != i2 && item.Label == item2.Label {
				return errors.Join(ErrDuplicatedItemLabel, errors.New(item.Label))
			}

			if item.Noun.Label == item2.Noun.Label && item.Adjective.Label == item2.Adjective.Label {
				return errors.Join(ErrDuplicatedNounAdj, errors.New(item.Label))
			}
		}
	}

	return nil
}
