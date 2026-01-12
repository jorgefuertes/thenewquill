package word

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
)

type query struct {
	db      *database.DB
	filters []database.Filter
}

func (s *Service) Get() *query {
	return &query{
		db: s.db,
		filters: []database.Filter{
			database.FilterByKind(kind.Word),
		},
	}
}

func (q *query) WithID(id uint32) *query {
	q.filters = append(q.filters, database.FilterByID(id))

	return q
}

func (q *query) WithNoID(id uint32) *query {
	q.filters = append(q.filters, database.NewFilter("ID", database.NotEqual, id))

	return q
}

func (q *query) WithLabel(label string) *query {
	q.filters = append(q.filters, database.FilterByLabel(label))

	return q
}

func (q *query) WithLabelID(id uint32) *query {
	q.filters = append(q.filters, database.FilterByLabelID(id))

	return q
}

func (q *query) WithSynonym(syn string) *query {
	q.filters = append(q.filters, database.NewFilter("Synonyms", database.Contains, syn))

	return q
}

func (q *query) WithType(t WordType) *query {
	q.filters = append(q.filters, database.NewFilter("Type", database.Equal, t))

	return q
}

func (q *query) Exists() bool {
	return q.db.Query(q.filters...).Exists()
}

func (q *query) First() (*Word, error) {
	w := &Word{}

	err := q.db.Query(q.filters...).First(w)

	return w, err
}

func (q *query) All() []*Word {
	var words []*Word

	cur := q.db.Query(q.filters...)
	w := &Word{}
	for cur.Next(w) {
		words = append(words, w)
		w = &Word{}
	}

	return words
}

func (q *query) Count() int {
	return q.db.Query(q.filters...).Count()
}

// GetAnyWith tries to find a Word by its label or synonym for any of the provided word types.
func (s *Service) GetAnyWith(labelOrSynonym string, wordTypes ...WordType) (*Word, error) {
	for _, t := range wordTypes {
		if w, err := s.Get().WithLabel(labelOrSynonym).WithType(t).First(); err == nil {
			return w, nil
		}

		if w, err := s.Get().WithSynonym(labelOrSynonym).WithType(t).First(); err == nil {
			return w, nil
		}
	}

	return nil, ErrWordNotFound
}
