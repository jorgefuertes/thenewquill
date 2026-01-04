package character

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
			database.FilterByKind(kind.Character),
		},
	}
}

func (q *query) WithHuman(h bool) *query {
	q.filters = append(q.filters, database.NewFilter("Human", database.Equal, h))

	return q
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

func (q *query) Exists() bool {
	return q.db.Query(q.filters...).Exists()
}

func (q *query) First() (*Character, error) {
	p := &Character{}

	err := q.db.Query(q.filters...).First(p)

	return p, err
}

func (q *query) Count() int {
	return q.db.Query(q.filters...).Count()
}
