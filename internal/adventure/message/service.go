package message

import "github.com/jorgefuertes/thenewquill/internal/adventure/db"

type Service struct {
	db *db.DB
}

func NewService(d *db.DB) *Service {
	return &Service{db: d}
}

func (s *Service) Create(m Message) error {
	if err := s.db.Append(m); err != nil {
		return err
	}

	return nil
}

func (s *Service) Update(m Message) error {
	return s.db.Update(m)
}

func (s *Service) Get(id db.ID) (Message, error) {
	i := Message{}
	err := s.db.Get(id, db.Messages, &i)

	return i, err
}

func (s *Service) All() []Message {
	words := make([]Message, 0)

	q := s.db.Query(db.Messages)
	var word Message
	for q.Next(&word) {
		words = append(words, word)
	}

	return words
}

func (s *Service) FindByLabel(labelName string) (Message, error) {
	label, err := s.db.GetLabelByName(labelName)
	if err != nil {
		return Message{}, err
	}

	return s.Get(label.ID)
}

func (s *Service) Count() int {
	return s.db.CountByKind(db.Messages)
}
