package db

import (
	"errors"
	"fmt"
	"math"
	"regexp"
)

const (
	Limit              = math.MaxInt32
	MinMeaningfulID ID = 4
)

type ID uint32

func (id ID) String() string {
	return fmt.Sprintf("%d", id)
}

func (id ID) UInt32() uint32 {
	return uint32(id)
}

type Label struct {
	ID   ID
	Name string
}

var (
	UndefinedLabel  = Label{ID: 0, Name: "undefined"}
	ConfigLabel     = Label{ID: 1, Name: "adv-config"}
	UnderscoreLabel = Label{ID: 2, Name: "_"}
	WildcardLabel   = Label{ID: 3, Name: "*"}
)

func (d *DB) AddLabel(name string) (Label, error) {
	if d.ExistsLabelName(name) {
		return d.GetLabelByName(name)
	}

	if !IsValidLabelName(name) {
		return Label{}, errors.Join(ErrInvalidLabelName, errors.New(name))
	}

	d.mut.Lock()
	defer d.mut.Unlock()

	if d.nextID >= Limit {
		return Label{}, ErrLimitReached
	}

	label := Label{ID: d.nextID, Name: name}
	d.Labels = append(d.Labels, label)
	d.nextID++

	return label, nil
}

func (d *DB) GetLabel(id ID) (Label, error) {
	d.mut.Lock()
	defer d.mut.Unlock()

	for _, l := range d.Labels {
		if l.ID == id {
			return l, nil
		}
	}

	return Label{ID: id, Name: "not-found"}, ErrNotFound
}

func (d *DB) GetLabelName(id ID) string {
	l, err := d.GetLabel(id)
	if err != nil {
		return UndefinedLabel.Name
	}

	return l.Name
}

func (d *DB) GetLabelByName(name string) (Label, error) {
	d.mut.Lock()
	defer d.mut.Unlock()

	for _, l := range d.Labels {
		if l.Name == name {
			return l, nil
		}
	}

	return Label{}, ErrNotFound
}

func (d *DB) ExistsLabelName(name string) bool {
	d.mut.Lock()
	defer d.mut.Unlock()

	for _, l := range d.Labels {
		if l.Name == name {
			return true
		}
	}

	return false
}

func IsValidLabelName(name string) bool {
	return regexp.MustCompile(`^[\d\p{L}\-_]{1,25}$`).MatchString(name)
}
