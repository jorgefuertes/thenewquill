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
	c := &Character{}
	err := s.db.GetByLabel(label, c)

	return c, err
}

func (s *Service) GetHuman() (*Character, error) {
	chars := s.db.Query(database.FilterByKind(kind.Character), database.NewFilter("Human", database.Equal, true))
	defer chars.Close()

	c := &Character{}
	if err := chars.First(c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *Service) HasHuman() bool {
	_, err := s.GetHuman()

	return err == nil
}

func (s *Service) Count() int {
	return s.db.CountRecordsByKind(kind.Character)
}
