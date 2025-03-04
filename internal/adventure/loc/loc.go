package loc

import "thenewquill/internal/adventure/words"

const Undefined = `undefined`

type Location struct {
	Label       string
	Title       string
	Description string
	Conns       []Connection
}

func New(label, title, desc string) *Location {
	return &Location{
		Label:       label,
		Title:       title,
		Description: desc,
		Conns:       make([]Connection, 0),
	}
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
