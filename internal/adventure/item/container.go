package item

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

func (s *Service) Weight(i Item) int {
	if !i.Container {
		return i.Weight
	}

	w := i.Weight

	items := s.db.Query(db.Items, db.Filter("At", i.ID), db.Filter("Container", true))

	var item Item
	for items.Next(&item) {
		w += s.Weight(item)
	}

	return w
}

func (s *Service) Move(i *Item, to db.Storeable) error {
	if s.IsContained(*i) {
		return ErrItemAlreadyContained
	}

	switch to.GetKind() {
	case db.Items:
		container, ok := to.(Item)
		if !ok {
			return ErrCannotAssertIntoItem
		}

		if s.Weight(*i)+s.Weight(container) > container.MaxWeight {
			return ErrContainerCantCarrySoMuch
		}

		i.At = to.GetID()
	case db.Locations, db.Characters:
		i.At = to.GetID()
	default:
		return ErrInvalidTo
	}

	return s.Update(*i)
}

func (s *Service) GetItemContainer(item Item) (Item, error) {
	return s.Get(item.At)
}

// IsContained returns true if the given item is contained in any container
func (s *Service) IsContained(item Item) bool {
	if item.At == db.UndefinedLabel.ID {
		return false
	}

	_, err := s.GetItemContainer(item)

	return err != nil
}

func (s *Service) Contents(id db.ID) []Item {
	items := make([]Item, 0)

	var item Item
	q := s.db.Query(db.Items, db.Filter("At", id))
	for q.Next(&item) {
		items = append(items, item)
	}

	return items
}
