package loc

import (
	"errors"
	"fmt"

	"thenewquill/internal/adventure/voc"
)

type Connection struct {
	Word *voc.Word
	To   *Location
}

type Location struct {
	Label       string
	Title       string
	Description string
	Connections []Connection
}

func NewLocation(label, title, description string) Location {
	return Location{
		Label:       label,
		Title:       title,
		Description: description,
		Connections: make([]Connection, 0),
	}
}

type Locations []Location

func New() Locations {
	return Locations{}
}

func (locs *Locations) Add(label, title, description string, connections []Connection) error {
	if locs.Exists(label) {
		return ErrDuplicatedLocation
	}

	*locs = append(*locs, Location{Label: label, Title: title, Description: description, Connections: connections})

	return nil
}

func (locs Locations) Get(label string) *Location {
	for _, l := range locs {
		if l.Label == label {
			return &l
		}
	}

	return nil
}

func (locs Locations) Exists(label string) bool {
	return locs.Get(label) != nil
}

func (locs Locations) getID(label string) (int, error) {
	for i, l := range locs {
		if l.Label == label {
			return i, nil
		}
	}

	return 0, ErrLocationNotFound
}

func (locs Locations) AddConnection(fromLabel string, word *voc.Word, toLabel string) error {
	to := locs.Get(toLabel)
	if to == nil {
		return errors.Join(fmt.Errorf("cannot find location '%s'", toLabel), ErrLocationNotFound)
	}

	id, err := locs.getID(fromLabel)
	if err != nil {
		return errors.Join(fmt.Errorf("cannot find id for '%s'", fromLabel), ErrLocationNotFound)
	}

	locs[id].Connections = append(locs[id].Connections, Connection{Word: word, To: to})

	return nil
}

func (l Location) Exits() []voc.Word {
	exits := make([]voc.Word, 0)

	for _, c := range l.Connections {
		exits = append(exits, *c.Word)
	}

	return exits
}

func (l Location) Go(word *voc.Word) *Location {
	for _, c := range l.Connections {
		if c.Word.IsEqual(word) {
			return c.To
		}
	}

	return nil
}
