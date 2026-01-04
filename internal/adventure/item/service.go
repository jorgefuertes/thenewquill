package item

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
)

type Service struct {
	db *database.DB
}

func NewService(d *database.DB) *Service {
	return &Service{db: d}
}

func (s *Service) Create(i *Item) (uint32, error) {
	return s.db.Create(i)
}

func (s *Service) Update(i *Item) error {
	return s.db.Update(i)
}

func (s *Service) Count() int {
	return s.db.CountRecordsByKind(kind.Item)
}

func (s *Service) SetCreated(i *Item, created bool) error {
	i.Created = created

	return s.Update(i)
}
