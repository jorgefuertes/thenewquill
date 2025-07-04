package item

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

type Service struct {
	db *db.DB
}

func NewService(d *db.DB) *Service {
	return &Service{db: d}
}

func (s *Service) Create(i Item) error {
	if err := s.db.Append(i); err != nil {
		return err
	}

	return nil
}

func (s *Service) Update(i Item) error {
	return s.db.Update(i)
}

func (s *Service) Get(id db.ID) (Item, error) {
	i := Item{}
	err := s.db.Get(id, &i)

	return i, err
}

func (s *Service) All() []Item {
	items := make([]Item, 0)

	q := s.db.Query(db.FilterByKind(db.Items))
	var item Item
	for q.Next(&item) {
		items = append(items, item)
	}

	return items
}

func (s *Service) SetCreated(i *Item, created bool) error {
	i.Created = created

	return s.Update(*i)
}

func (s *Service) IsAt(i Item, at db.Storeable) bool {
	return i.At == at.GetID()
}

func (s *Service) FindByLabel(labelName string) (Item, error) {
	label, err := s.db.GetLabelByName(labelName)
	if err != nil {
		return Item{}, err
	}

	return s.Get(label.ID)
}

func (s *Service) Count() int {
	return s.db.CountByKind(db.Items)
}
