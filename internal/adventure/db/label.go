package db

import (
	"math"

	"github.com/jorgefuertes/thenewquill/internal/adventure/id"
	"github.com/jorgefuertes/thenewquill/internal/adventure/label"
)

func (d *DB) AddLabel(name string) (label.Label, error) {
	if d.ExistsLabelName(name) {
		return d.GetLabelByName(name)
	}

	l := label.Label{ID: d.nextID, Name: name}
	if err := l.Validate(false); err != nil {
		return l, err
	}

	d.lock()
	defer d.unlock()

	d.Labels = append(d.Labels, l)

	if d.nextID == math.MaxUint32 {
		return label.Label{}, ErrLimitReached
	}

	d.nextID++

	return l, nil
}

func (d *DB) GetLabel(id id.ID) (label.Label, error) {
	d.lock()
	defer d.unlock()

	for _, l := range d.Labels {
		if l.ID == id {
			return l, nil
		}
	}

	return label.Label{ID: id, Name: "not-found"}, ErrNotFound
}

func (d *DB) GetLabelName(id id.ID) string {
	l, err := d.GetLabel(id)
	if err != nil {
		return label.Undefined.Name
	}

	return l.Name
}

func (d *DB) GetLabelByName(name string) (label.Label, error) {
	d.lock()
	defer d.unlock()

	for _, l := range d.Labels {
		if l.Name == name {
			return l, nil
		}
	}

	return label.Label{}, ErrNotFound
}

func (d *DB) ExistsLabelName(name string) bool {
	d.lock()
	defer d.unlock()

	for _, l := range d.Labels {
		if l.Name == name {
			return true
		}
	}

	return false
}
