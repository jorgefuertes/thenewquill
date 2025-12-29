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

type query struct {
	db    *database.DB
	id    uint32
	label string
	syn   string
	t     WordType
}

func (s *Service) Get() *query {
	return &query{db: s.db}
}

func (q *query) WithID(id uint32) *query {
	q.id = id

	return q
}

func (q *query) WithLabel(label string) *query {
	q.label = label

	return q
}

func (q *query) WithSynonym(syn string) *query {
	q.syn = syn

	return q
}

func (q *query) WithType(t WordType) *query {
	q.t = t

	return q
}

func (q *query) First() (*Word, error) {
	w := &Word{}

	filters := []database.Filter{database.FilterByKind(kind.Word)}

	if q.id != 0 {
		filters = append(filters, database.FilterByID(q.id))
	}

	if q.t != None {
		filters = append(filters, database.NewFilter("Type", database.Equal, q.t))
	}

	if q.label != "" {
		filters = append(filters, database.FilterByLabel(q.label))
	}

	if q.syn != "" {
		filters = append(filters, database.NewFilter("Synonyms", database.Contains, q.syn))
	}

	err := q.db.Query(filters...).First(w)

	return w, err
}
