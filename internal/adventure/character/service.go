package character

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

type Service struct {
	db *db.DB
}

func NewService(d *db.DB) *Service {
	return &Service{db: d}
}

func (s *Service) Create(c Character) error {
	if err := s.db.Append(c); err != nil {
		return err
	}

	return nil
}

func (s *Service) Update(c Character) error {
	return s.db.Update(c)
}

func (s *Service) Get(id db.ID) (Character, error) {
	c := Character{}
	err := s.db.Get(id, db.Characters, &c)

	return c, err
}

func (s *Service) GetHuman() (Character, error) {
	chars := s.db.Query(db.Characters, db.Filter("Human", true))
	defer chars.Close()

	var c Character
	if !chars.Next(&c) {
		return Character{}, ErrNoHuman
	}

	return c, nil
}

func (s *Service) HasHuman() bool {
	_, err := s.GetHuman()

	return err == nil
}

func (s *Service) FindByLabel(labelName string) (Character, error) {
	label, err := s.db.GetLabelByName(labelName)
	if err != nil {
		return Character{}, err
	}

	return s.Get(label.ID)
}

func (s *Service) Count() int {
	return s.db.CountByKind(db.Characters)
}
