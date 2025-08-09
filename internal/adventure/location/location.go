package location

import (
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/util"
)

const Undefined = `undefined`

type Location struct {
	ID          db.ID
	Title       string
	Description string
	Conns       []Connection
}

var _ db.Storeable = &Location{}

func (l Location) Export() string {
	out := fmt.Sprintf("%d|%d|%s|%s",
		kind.KindOf(l).Byte(),
		l.ID,
		util.EscapeExportString(l.Title),
		util.EscapeExportString(l.Description),
	)

	if len(l.Conns) == 0 {
		return out + "\n"
	}

	out += "|"

	for i := 0; i < len(l.Conns); i++ {
		out += fmt.Sprintf("%d:%d", l.Conns[i].WordID, l.Conns[i].LocationID)
		if i != len(l.Conns)-1 {
			out += ","
		}
	}

	return out + "\n"
}

func New(title, desc string) Location {
	return Location{
		ID:          db.UndefinedLabel.ID,
		Title:       title,
		Description: desc,
		Conns:       make([]Connection, 0),
	}
}

func (l Location) SetID(id db.ID) db.Storeable {
	l.ID = id

	return l
}

func (l Location) GetID() db.ID {
	return l.ID
}

func (l *Location) connIndex(wordID db.ID) int {
	for i, c := range l.Conns {
		if c.WordID == wordID {
			return i
		}
	}

	return -1
}

func (l *Location) SetConn(wordID, locationID db.ID) {
	idx := l.connIndex(wordID)
	if idx != -1 {
		l.Conns[idx].LocationID = locationID

		return
	}

	l.Conns = append(l.Conns, Connection{WordID: wordID, LocationID: locationID})
}

func (l *Location) GetConn(wordID db.ID) db.ID {
	idx := l.connIndex(wordID)
	if idx != -1 {
		return l.Conns[idx].LocationID
	}

	return db.UndefinedLabel.ID
}

func (l *Location) HasConn(wordID db.ID) bool {
	return l.connIndex(wordID) != -1
}
