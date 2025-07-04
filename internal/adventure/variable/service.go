package variable

import (
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

type Service struct {
	db *db.DB
}

func NewService(d *db.DB) *Service {
	return &Service{db: d}
}

func (s *Service) Set(id db.ID, value any) error {
	v := Variable{ID: id, Value: value}
	if s.db.Exists(id) {
		return s.db.Update(v)
	}

	return s.db.Append(v)
}

func (s *Service) Get(id db.ID) (Variable, error) {
	var v Variable
	err := s.db.Get(id, &v)

	return v, err
}

func (s *Service) All() []Variable {
	vars := make([]Variable, 0)

	q := s.db.Query(db.FilterByKind(db.Variables))
	var varr Variable
	for q.Next(&varr) {
		vars = append(vars, varr)
	}

	return vars
}

func (s *Service) FindByLabel(paths ...string) (Variable, error) {
	label, err := s.db.GetLabelByName(strings.Join(paths, "."))
	if err != nil {
		return Variable{}, err
	}

	return s.Get(label.ID)
}

func (s *Service) Count() int {
	return s.db.CountByKind(db.Variables)
}
