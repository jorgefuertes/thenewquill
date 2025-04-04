package loc

import (
	"strings"

	"thenewquill/internal/adventure/vars"
	"thenewquill/internal/adventure/words"
)

const Undefined = `undefined`

type Location struct {
	Label       string
	Title       string
	Description string
	Conns       []Connection
	Vars        vars.Store
}

func New(label, title, desc string) *Location {
	return &Location{
		Label:       label,
		Title:       title,
		Description: desc,
		Conns:       make([]Connection, 0),
		Vars:        vars.NewStore(),
	}
}

func (l *Location) GetLabel() string {
	return l.Label
}

func (l *Location) connIndex(word *words.Word) int {
	for i, c := range l.Conns {
		if c.Word.Is(word.Label) {
			return i
		}
	}

	return -1
}

func (l *Location) SetConn(word *words.Word, to *Location) {
	idx := l.connIndex(word)
	if idx != -1 {
		l.Conns[idx].To = to

		return
	}

	l.Conns = append(l.Conns, Connection{Word: word, To: to})
}

func (l *Location) GetConn(word *words.Word) *Location {
	idx := l.connIndex(word)
	if idx != -1 {
		return l.Conns[idx].To
	}

	return nil
}

func (l *Location) toLower() {
	l.Label = strings.ToLower(l.Label)
}
