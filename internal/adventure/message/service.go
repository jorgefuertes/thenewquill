package message

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

func (s *Service) Create(msg *Message) (uint32, error) {
	return s.db.Create(msg)
}

func (s *Service) Update(msg *Message) error {
	return s.db.Update(msg)
}

func (s *Service) Get(id uint32) (*Message, error) {
	msg := &Message{}
	err := s.db.Get(id, &msg)

	return msg, err
}

func (s *Service) GetByLabel(label string) (*Message, error) {
	msg := &Message{}
	err := s.db.GetByLabel(label, msg)

	return msg, err
}

func (s *Service) Count() int {
	return s.db.CountRecordsByKind(kind.Message)
}

func (s *Service) GetHuman() (*Message, error) {
	chars := s.db.Query(database.FilterByKind(kind.Message), database.NewFilter("Human", database.Equal, true))
	defer chars.Close()

	var loc *Message
	err := chars.First(loc)

	return loc, err
}

func (s *Service) HasHuman() bool {
	_, err := s.GetHuman()

	return err == nil
}
