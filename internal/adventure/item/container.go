package item

import (
	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

func (s *Service) Weight(i Item) int {
	if !i.Container {
		return i.Weight
	}

	w := i.Weight

	items := s.db.Query(
		database.FilterByKind(kind.Item),
		database.Filter("At", database.Equal, i.ID),
		database.Filter("Container", database.Equal, true),
	)

	var item Item
	for items.Next(&item) {
		w += s.Weight(item)
	}

	return w
}

func (s *Service) Move(i *Item, to adapter.Storeable) error {
	if s.IsContained(*i) {
		return ErrItemAlreadyContained
	}

	switch kind.KindOf(to) {
	case kind.Item:
		container, ok := to.(*Item)
		if !ok {
			return ErrCannotAssertIntoItem
		}

		if s.Weight(*i)+s.Weight(*container) > container.MaxWeight {
			return ErrContainerCantCarrySoMuch
		}

		i.At = to.GetID()
	case kind.Location, kind.Character:
		i.At = to.GetID()
	default:
		return ErrInvalidTo
	}

	return s.Update(i)
}

func (s *Service) GetItemContainer(item Item) (*Item, error) {
	return s.Get(item.At)
}

// IsContained returns true if the given item is contained in any container
func (s *Service) IsContained(item Item) bool {
	if item.At == primitive.UndefinedID {
		return false
	}

	_, err := s.GetItemContainer(item)

	return err != nil
}

func (s *Service) Contents(id primitive.ID) []Item {
	items := make([]Item, 0)

	var item Item
	q := s.db.Query(database.FilterByKind(kind.Item), database.Filter("At", database.Equal, id))
	for q.Next(&item) {
		items = append(items, item)
	}

	return items
}
