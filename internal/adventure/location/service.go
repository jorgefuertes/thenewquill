package location

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(loc *Location) (uint32, error) {
	return s.db.Create(loc)
}

func (s *Service) Update(loc *Location) error {
	return s.db.Update(loc)
}

func (s *Service) Count() int {
	return s.db.CountRecordsByKind(kind.Location)
}
