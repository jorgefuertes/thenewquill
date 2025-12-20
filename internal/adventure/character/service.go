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

func (s *Service) DB() *database.DB {
	return s.db
}

func (s *Service) GetHuman() (*Character, error) {
	chars := s.db.Query(database.FilterByKind(kind.Character), database.Filter("Human", database.Equal, true))
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
	return s.db.Count(database.FilterByKind(kind.Character))
}
