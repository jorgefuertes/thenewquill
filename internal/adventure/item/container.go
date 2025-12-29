package item

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
)

func (s *Service) TotalWeight(i Item) int {
	if !i.Container {
		return i.Weight
	}

	w := i.Weight

	items := s.db.Query(
		database.FilterByKind(kind.Item),
		database.NewFilter("At", database.Equal, i.ID),
		database.NewFilter("Container", database.Equal, true),
	)

	var item Item
	for items.Next(&item) {
		w += s.TotalWeight(item)
	}

	return w
}

func (s *Service) PutInto(i *Item, in Item) error {
	if s.IsContained(*i) {
		return ErrItemAlreadyContained
	}

	if s.TotalWeight(*i)+s.TotalWeight(in) > in.MaxWeight {
		return ErrContainerCantCarrySoMuch
	}

	i.At = in.ID

	return s.Update(i)
}

func (s *Service) GetItemContainer(item Item) (*Item, error) {
	return s.Get(item.At)
}

// IsContained returns true if the given item is contained in any container
func (s *Service) IsContained(item Item) bool {
	if item.At == 0 {
		return false
	}

	_, err := s.GetItemContainer(item)

	return err != nil
}

func (s *Service) Contents(id uint32) []Item {
	items := make([]Item, 0)

	var item Item
	q := s.db.Query(database.FilterByKind(kind.Item), database.NewFilter("At", database.Equal, id))
	for q.Next(&item) {
		items = append(items, item)
	}

	return items
}
