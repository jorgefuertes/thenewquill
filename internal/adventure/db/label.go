package db

import (
	"errors"
	"fmt"
	"math"
	"regexp"
)

const (
	Limit                  = math.MaxUint32
	MinMeaningfulID  ID    = 4
	AllowSpecial     Allow = true
	DontAllowSpecial Allow = false
	AllowDot         Allow = true
	DontAllowDot     Allow = false
)

type ID uint32

func (id ID) String() string {
	return fmt.Sprintf("%d", id)
}

func (id ID) UInt32() uint32 {
	return uint32(id)
}

func (id ID) IsDefined() bool {
	return id != UndefinedLabel.ID
}

// Validate checks if the ID is valid, allowSpecial allows _ and * to be valid IDs
func (id ID) Validate(allowSpecial Allow) error {
	if id == UndefinedLabel.ID {
		return ErrUndefinedLabel
	}

	if id < MinMeaningfulID && !allowSpecial {
		return ErrInvalidLabelID
	}

	if id >= Limit {
		return ErrInvalidLabelID
	}

	return nil
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

func (d *DB) AddLabel(name string, allowDot Allow) (Label, error) {
	if d.ExistsLabelName(name) {
		return d.GetLabelByName(name)
	}

	if !IsValidLabelName(name, allowDot) {
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

func IsValidLabelName(name string, allowDot Allow) bool {
	if allowDot {
		return regexp.MustCompile(`^[\d\p{L}\-_\.]{1,25}$`).MatchString(name)
	}

	return regexp.MustCompile(`^[\d\p{L}\-_]{1,25}$`).MatchString(name)
}
