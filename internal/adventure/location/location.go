package location

import "github.com/jorgefuertes/thenewquill/internal/database/adapter"

const Undefined = `undefined`

type Location struct {
	ID          uint32
	LabelID     uint32 `valid:"required"`
	Title       string `valid:"required"`
	Description string `valid:"required"`
	Conns       []Connection
}

var _ adapter.Storeable = &Location{}

func New() *Location {
	return &Location{Conns: make([]Connection, 0)}
}

func (l *Location) SetID(id uint32) {
	l.ID = id
}

func (l *Location) GetID() uint32 {
	return l.ID
}

func (l *Location) SetLabelID(id uint32) {
	l.LabelID = id
}

func (l *Location) GetLabelID() uint32 {
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

	return 0
}

func (l *Location) HasConn(wordID uint32) bool {
	return l.connIndex(wordID) != -1
}
