package character

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

func (s *Service) GetByLabel(label string) (*Character, error) {
	return s.Get().WithLabel(label).First()
}

func (s *Service) GetHuman() (*Character, error) {
	return s.Get().WithHuman(true).First()
}

func (s *Service) HasHuman() bool {
	return s.Get().WithHuman(true).Exists()
}

func (s *Service) Count() int {
	return s.db.CountRecordsByKind(kind.Character)
}
