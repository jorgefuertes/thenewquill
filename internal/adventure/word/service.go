package word

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/database"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

type Service struct {
	db *database.DB
}

func NewService(d *database.DB) *Service {
	return &Service{db: d}
}

func (s *Service) Create(w *Word) (primitive.ID, error) {
	return s.db.Create(w)
}

func (s *Service) Update(w *Word) error {
	return s.db.Update(w)
}

func (s *Service) Get(id primitive.ID) (*Word, error) {
	w := &Word{}
	err := s.db.Get(id, &w)

	return w, err
}

func (s *Service) GetByLabel(label primitive.Label) (*Word, error) {
	w := &Word{}
	err := s.db.GetByLabel(label, w)

	return w, err
}

func (s *Service) First(label primitive.Label) (*Word, error) {
	cursor := s.db.Query(database.FilterByKind(kind.Word))
	defer cursor.Close()

	labelID, err := s.db.GetLabelID(label)
	if err != nil {
		return nil, err
	}

	w := &Word{}
	for cursor.Next(w) {
		if w.LabelID == labelID {
			return w, nil
		}

		for _, syn := range w.Synonyms {
			if syn == label.String() {
				return w, nil
			}
		}
	}

	return nil, database.ErrNotFound
}

func (s *Service) Count() int {
	return s.db.Count(database.FilterByKind(kind.Word))
}
