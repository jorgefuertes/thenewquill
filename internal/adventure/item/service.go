package item

import (
	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

type Service struct {
	db *database.DB
}

func NewService(d *database.DB) *Service {
	return &Service{db: d}
}

func (s *Service) DB() *database.DB {
	return s.db
}

func (s *Service) Create(i *Item) (primitive.ID, error) {
	return s.db.Create(i)
}

func (s *Service) Update(i *Item) error {
	return s.db.Update(i)
}

func (s *Service) Get(id primitive.ID) (*Item, error) {
	i := &Item{}
	err := s.db.Get(id, &i)

	return i, err
}

func (s *Service) GetByLabel(labelOrString any) (*Item, error) {
	label, err := primitive.LabelFromLabelOrString(labelOrString)
	if err != nil {
		return nil, err
	}

	i := &Item{}
	err = s.db.GetByLabel(label, i)

	return i, err
}

func (s *Service) Count() int {
	return s.db.Count(database.FilterByKind(kind.Item))
}

func (s *Service) SetCreated(i *Item, created bool) error {
	i.Created = created

	return s.Update(i)
}

func (s *Service) IsAt(i Item, at adapter.Storeable) bool {
	return i.At == at.GetID()
}
