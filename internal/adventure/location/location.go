package location

import (
	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/database/primitive"
)

const Undefined = `undefined`

type Location struct {
	ID          uint32
	LabelID     uint32
	Title       string
	Description string
	Conns       []Connection
}

var _ adapter.Storeable = &Location{}

func New(title, desc string) *Location {
	return &Location{
		ID:          primitive.UndefinedID,
		LabelID:     primitive.UndefinedID,
		Title:       title,
		Description: desc,
		Conns:       make([]Connection, 0),
	}
}

func (l *Location) SetID(id uint32) {
	l.ID = id
}

func (l Location) GetID() uint32 {
	return l.ID
}

func (l *Location) SetLabelID(id uint32) {
	l.LabelID = id
}

func (l Location) GetLabelID() uint32 {
	return l.LabelID
}

func (l *Location) connIndex(wordID uint32) int {
	for i, c := range l.Conns {
		if c.WordID == wordID {
			return i
		}
	}

	return -1
}

func (l *Location) SetConn(wordID, locationID uint32) {
	idx := l.connIndex(wordID)
	if idx != -1 {
		l.Conns[idx].LocationID = locationID

		return
	}

	l.Conns = append(l.Conns, Connection{WordID: wordID, LocationID: locationID})
}

func (l *Location) GetConn(wordID uint32) uint32 {
	idx := l.connIndex(wordID)
	if idx != -1 {
		return l.Conns[idx].LocationID
	}

	return primitive.UndefinedID
}

func (l *Location) HasConn(wordID uint32) bool {
	return l.connIndex(wordID) != -1
}
