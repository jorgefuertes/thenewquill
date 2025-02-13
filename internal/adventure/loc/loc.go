package loc

import "thenewquill/internal/adventure/voc"

type Location struct {
	Label       string
	Title       string
	Description string
	Conns       []Connection
}

func NewLocation(label, title, desc string) *Location {
	return &Location{
		Label:       label,
		Title:       title,
		Description: desc,
		Conns:       make([]Connection, 0),
	}
}

func (l *Location) connIndex(word *voc.Word) int {
	for i, c := range l.Conns {
		if c.Word == word {
			return i
		}
	}

	return -1
}

func (l *Location) SetConn(word *voc.Word, to *Location) {
	idx := l.connIndex(word)
	if idx != -1 {
		l.Conns[idx].To = to

		return
	}

	l.Conns = append(l.Conns, Connection{Word: word, To: to})
}

func (l *Location) GetConn(word *voc.Word) *Location {
	idx := l.connIndex(word)
	if idx != -1 {
		return l.Conns[idx].To
	}

	return nil
}
