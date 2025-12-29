package variable

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/internal/util"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Set(id uint32, value any) (uint32, error) {
	v := &Variable{}

	err := s.db.Get(id, v)
	if err != nil {
		return 0, err
	}

	v.SetValue(value)

	return v.ID, s.db.Update(v)
}

func (s *Service) SetByLabel(label string, value any) (uint32, error) {
	labelID, err := s.db.CreateLabel(label)
	if err != nil {
		return 0, err
	}

	v := &Variable{}

	err = s.db.GetByLabel(label, v)
	if err == nil {
		v.SetValue(value)

		return v.ID, s.db.Update(v)
	}

	v.LabelID = labelID
	v.SetValue(value)

	return s.db.Create(v)
}

func (s *Service) Count() int {
	return s.db.CountRecordsByKind(kind.Variable)
}

type query struct {
	db      *database.DB
	id      uint32
	label   string
	labelID uint32
	value   string
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

func (q *query) WithValue(value any) *query {
	q.value = util.ValueToString(value)

	return q
}

func (q *query) First() (*Variable, error) {
	w := &Variable{}

	filters := []database.Filter{database.FilterByKind(kind.Variable)}

	if q.id != 0 {
		filters = append(filters, database.FilterByID(q.id))
	}

	if q.label != "" {
		filters = append(filters, database.FilterByLabel(q.label))
	}

	if q.labelID != 0 {
		filters = append(filters, database.NewFilter("LabelID", database.Equal, q.labelID))
	}

	if q.value != "" {
		filters = append(filters, database.NewFilter("Value", database.Equal, q.value))
	}

	err := q.db.Query(filters...).First(w)

	return w, err
}
