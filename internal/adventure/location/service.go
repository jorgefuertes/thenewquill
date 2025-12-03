package location

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

func (s *Service) Create(loc *Location) (primitive.ID, error) {
	return s.db.Create(loc)
}

func (s *Service) Update(loc *Location) error {
	return s.db.Update(loc)
}

func (s *Service) Get(id primitive.ID) (*Location, error) {
	loc := &Location{}
	err := s.db.Get(id, &loc)

	return loc, err
}

func (s *Service) GetByLabel(labelOrString any) (*Location, error) {
	label, err := primitive.LabelFromLabelOrString(labelOrString)
	if err != nil {
		return nil, err
	}

	loc := &Location{}
	err = s.db.GetByLabel(label, loc)

	return loc, err
}

func (s *Service) Count() int {
	return s.db.Count(database.FilterByKind(kind.Location))
}
