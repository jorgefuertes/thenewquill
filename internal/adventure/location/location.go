package location

import (
	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
)

const Undefined = `undefined`

type Location struct {
	ID          primitive.ID
	LabelID     primitive.ID
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

func (l *Location) SetID(id primitive.ID) {
	l.ID = id
}

func (l Location) GetID() primitive.ID {
	return l.ID
}

func (l *Location) SetLabelID(id primitive.ID) {
	l.LabelID = id
}

func (l Location) GetLabelID() primitive.ID {
	return l.LabelID
}

func (l *Location) connIndex(wordID primitive.ID) int {
	for i, c := range l.Conns {
		if c.WordID == wordID {
			return i
		}
	}

	return -1
}

func (l *Location) SetConn(wordID, locationID primitive.ID) {
	idx := l.connIndex(wordID)
	if idx != -1 {
		l.Conns[idx].LocationID = locationID

		return
	}

	l.Conns = append(l.Conns, Connection{WordID: wordID, LocationID: locationID})
}

func (l *Location) GetConn(wordID primitive.ID) primitive.ID {
	idx := l.connIndex(wordID)
	if idx != -1 {
		return l.Conns[idx].LocationID
	}

	return primitive.UndefinedID
}

func (l *Location) HasConn(wordID primitive.ID) bool {
	return l.connIndex(wordID) != -1
}
