package message

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

func (s *Service) Create(msg *Message) (uint32, error) {
	return s.db.Create(msg)
}

func (s *Service) Update(msg *Message) error {
	return s.db.Update(msg)
}

func (s *Service) Count() int {
	return s.db.CountRecordsByKind(kind.Message)
}
