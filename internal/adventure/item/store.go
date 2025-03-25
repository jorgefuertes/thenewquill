package item

import (
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
	if label == "" {
		return -1
	}

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

func (s *Store) GetAll() []*Item {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.items
}

func (s *Store) Exists(label string) bool {
	return s.getIndex(label) > -1
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

// Set a new empty item
func (s *Store) New(label string) (*Item, error) {
	i := New(label, nil, nil)
	if err := s.Set(i); err != nil {
		return nil, err
	}

	return i, nil
}

// Set or replace an item
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
