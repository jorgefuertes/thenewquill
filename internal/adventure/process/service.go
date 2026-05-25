package process

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

func (s *Service) CreateTable(k TableKind) (uint32, error) {
	labelID, err := s.db.CreateLabel(k.String())
	if err != nil {
		return 0, err
	}

	t := &Table{
		LabelID: labelID,
		Kind:    k,
		Procs:   make([]Process, 0),
	}

	return s.db.Create(t)
}

func (s *Service) Update(t *Table) error {
	return s.db.Update(t)
}

func (s *Service) Count() int {
	return s.db.CountRecordsByKind(kind.Table)
}

func (s *Service) Exists(k TableKind) bool {
	return s.db.Query(database.FilterByKind(kind.Table), database.FilterByLabel(k.String())).Exists()
}
