package location

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(loc *Location) (uint32, error) {
	return s.db.Create(loc)
}

func (s *Service) Update(loc *Location) error {
	return s.db.Update(loc)
}

func (s *Service) Count() int {
	return s.db.CountRecordsByKind(kind.Location)
}

type query struct {
	db      *database.DB
	id      uint32
	label   string
	labelID uint32
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

func (q *query) WithLabelID(id uint32) *query {
	q.labelID = id

	return q
}

func (q *query) First() (*Location, error) {
	w := &Location{}

	filters := []database.Filter{database.FilterByKind(kind.Location)}

	if q.id != 0 {
		filters = append(filters, database.FilterByID(q.id))
	}

	if q.label != "" {
		filters = append(filters, database.FilterByLabel(q.label))
	}

	if q.labelID != 0 {
		filters = append(filters, database.NewFilter("LabelID", database.Equal, q.labelID))
	}

	err := q.db.Query(filters...).First(w)

	return w, err
}
