package variable

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/database"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(v *Variable) (primitive.ID, error) {
	return s.db.Create(v)
}

func (s *Service) CreateWithLabel(labelOrString any, value any) (primitive.ID, error) {
	labelID, err := s.db.CreateLabelIfNotExists(labelOrString, true)
	if err != nil {
		return primitive.UndefinedID, err
	}

	v := New(primitive.UndefinedID, labelID, value)

	return s.db.Create(v)
}

func (s *Service) Update(v *Variable) error {
	return s.db.Update(v)
}

func (s *Service) Get(id primitive.ID) (*Variable, error) {
	v := &Variable{}
	err := s.db.Get(id, v)

	return v, err
}

func (s *Service) GetByLabel(labelOrString any) (*Variable, error) {
	label, err := primitive.LabelFromLabelOrString(labelOrString)
	if err != nil {
		return nil, err
	}

	v := &Variable{}
	err = s.db.GetByLabel(label, v)

	return v, err
}

func (s *Service) Count() int {
	return s.db.Count(database.FilterByKind(kind.Variable))
}
