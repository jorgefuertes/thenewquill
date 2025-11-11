package location

import (
	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/id"
)

const Undefined = `undefined`

type Location struct {
	ID          id.ID
	Title       string
	Description string
	Conns       []Connection
}

var _ adapter.Storeable = &Location{}

func New(title, desc string) Location {
	return Location{
		ID:          id.Undefined,
		Title:       title,
		Description: desc,
		Conns:       make([]Connection, 0),
	}
}

func (l Location) SetID(id id.ID) adapter.Storeable {
	l.ID = id

	return l
}

func (l Location) GetID() id.ID {
	return l.ID
}

func (l *Location) connIndex(wordID id.ID) int {
	for i, c := range l.Conns {
		if c.WordID == wordID {
			return i
		}
	}

	return -1
}

func (l *Location) SetConn(wordID, locationID id.ID) {
	idx := l.connIndex(wordID)
	if idx != -1 {
		l.Conns[idx].LocationID = locationID

		return
	}

	l.Conns = append(l.Conns, Connection{WordID: wordID, LocationID: locationID})
}

func (l *Location) GetConn(wordID id.ID) id.ID {
	idx := l.connIndex(wordID)
	if idx != -1 {
		return l.Conns[idx].LocationID
	}

	return id.Undefined
}

func (l *Location) HasConn(wordID id.ID) bool {
	return l.connIndex(wordID) != -1
}
