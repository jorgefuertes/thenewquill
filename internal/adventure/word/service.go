package word

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
)

type Service struct {
	db *database.DB
}

func NewService(d *database.DB) *Service {
	s := &Service{db: d}

	labels := map[uint32]string{
		1: database.LabelAsterisk,
		2: database.LabelUnderscore,
	}

	for _, t := range WordTypes {
		for id, label := range labels {
			_, err := s.Create(&Word{
				LabelID:  id,
				Type:     t,
				Synonyms: []string{label},
			})
			if err != nil {
				panic(err)
			}
		}
	}

	return s
}

func (s *Service) Create(w *Word) (uint32, error) {
	return s.db.Create(w)
}

func (s *Service) Update(w *Word) error {
	return s.db.Update(w)
}

func (s *Service) Count() int {
	return s.db.CountRecordsByKind(kind.Word)
}
