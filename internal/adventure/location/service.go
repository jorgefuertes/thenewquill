package location

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

type Service struct {
	db *db.DB
}

func NewService(d *db.DB) *Service {
	return &Service{db: d}
}

func (s *Service) Create(l Location) error {
	if err := s.db.Append(l); err != nil {
		return err
	}

	return nil
}

func (s *Service) Update(l Location) error {
	return s.db.Update(l)
}

func (s *Service) Get(id db.ID) (Location, error) {
	i := Location{}
	err := s.db.Get(id, db.Locations, &i)

	return i, err
}

func (s *Service) All() []Location {
	locations := make([]Location, 0)

	q := s.db.Query(db.Locations)
	var location Location
	for q.Next(&location) {
		locations = append(locations, location)
	}

	return locations
}

func (s *Service) FindByLabel(labelName string) (Location, error) {
	label, err := s.db.GetLabelByName(labelName)
	if err != nil {
		return Location{}, err
	}

	return s.Get(label.ID)
}

func (s *Service) Count() int {
	return s.db.CountByKind(db.Locations)
}
