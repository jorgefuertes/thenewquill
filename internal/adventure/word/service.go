package word

import (
	"slices"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

type Service struct {
	db *db.DB
}

func NewService(d *db.DB) *Service {
	return &Service{db: d}
}

func (s *Service) Create(w Word) error {
	if err := s.db.Append(w); err != nil {
		return err
	}

	return nil
}

func (s *Service) Update(w Word) error {
	return s.db.Update(w)
}

func (s *Service) Get(id db.ID) (Word, error) {
	i := Word{}
	err := s.db.Get(id, db.Words, &i)

	return i, err
}

func (s *Service) All() []Word {
	words := make([]Word, 0)

	q := s.db.Query(db.Words)
	var word Word
	for q.Next(&word) {
		words = append(words, word)
	}

	return words
}

func (s *Service) FindByLabel(labelName string) (Word, error) {
	label, err := s.db.GetLabelByName(labelName)
	if err != nil {
		return Word{}, err
	}

	return s.Get(label.ID)
}

func (s *Service) First(syn string) (Word, error) {
	for _, w := range s.All() {
		if slices.Contains(w.Synonyms, syn) {
			return w, nil
		}
	}

	return Word{}, db.ErrNotFound
}

func (s *Service) Count() int {
	return s.db.CountByKind(db.Words)
}
