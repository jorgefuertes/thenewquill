package db

import (
	"reflect"
	"sync"

	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/id"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/adventure/label"
	"github.com/jorgefuertes/thenewquill/pkg/log"
)

const NotFound int = -1

type DB struct {
	nextID id.ID
	mut    *sync.Mutex
	Labels []label.Label
	Data   []adapter.Storeable
}

func (b *DB) lock() {
	b.mut.Lock()
}

func (b *DB) unlock() {
	b.mut.Unlock()
}

func New() *DB {
	return &DB{
		mut:    &sync.Mutex{},
		nextID: id.Min,
		Labels: []label.Label{label.Undefined, label.Config, label.Underscore, label.Wildcard},
		Data:   make([]adapter.Storeable, 0),
	}
}

func (d *DB) Reset() {
	d.lock()
	defer d.unlock()

	// reset labels
	d.Labels = []label.Label{label.Undefined, label.Config, label.Underscore, label.Wildcard}
	d.nextID = id.Min

	// reset entities
	d.Data = make([]adapter.Storeable, 0)
}

func (d *DB) Len() int {
	d.lock()
	defer d.unlock()

	return len(d.Data)
}

func (d *DB) indexOf(filters ...filter) int {
	for i, r := range d.Data {
		if matches(r, filters...) {
			return i
		}
	}

	return NotFound
}

func (d *DB) Exists(filters ...filter) bool {
	d.lock()
	defer d.unlock()

	return d.indexOf(filters...) != NotFound
}

func (d *DB) Create(labelName string, s adapter.Storeable) (id.ID, error) {
	if s.GetID() != label.Undefined.ID {
		return label.Undefined.ID, ErrCannotCreateWithDefinedID
	}

	if err := s.Validate(true); err != nil {
		return label.Undefined.ID, err
	}

	l, err := d.AddLabel(labelName)
	if err != nil {
		return label.Undefined.ID, err
	}

	s = s.SetID(l.ID)

	return l.ID, d.Append(s)
}

func (d *DB) Append(s adapter.Storeable) error {
	if kind.KindOf(s) == kind.None {
		return ErrKindCannotBeNone
	}

	if d.Exists(FilterByID(s.GetID()), FilterByKind(kind.KindOf(s))) {
		l, err := d.GetLabel(s.GetID())
		if err != nil {
			l = label.Label{Name: err.Error()}
		}

		k := kind.KindOf(s)
		log.Error("❗ duplicated record, %s %q", k, l.Name)

		old := s
		_ = d.Get(s.GetID(), &old)
		log.Error("❗ previous record kind %q: %v", kind.KindOf(old), old)

		return ErrDuplicatedRecord
	}

	if err := s.GetID().Validate(false); err != nil {
		return err
	}

	if err := s.Validate(false); err != nil {
		return err
	}

	d.lock()
	defer d.unlock()

	d.Data = append(d.Data, s)

	return nil
}

func (d *DB) Update(s adapter.Storeable) error {
	if !d.Exists(FilterByID(s.GetID()), FilterByKind(kind.KindOf(s))) {
		return ErrRecordNotFound
	}

	d.lock()
	defer d.unlock()

	idx := d.indexOf(FilterByID(s.GetID()), FilterByKind(kind.KindOf(s)))

	d.Data = append(d.Data[:idx], d.Data[idx+1:]...)
	d.Data = append(d.Data, s)

	return nil
}

func (d *DB) GetByLabel(labelName string, dst any) error {
	l, err := d.GetLabelByName(labelName)
	if err != nil {
		return err
	}

	return d.Get(l.ID, dst)
}

func (d *DB) Get(i id.ID, dst any) error {
	if i < id.Min {
		return ErrRecordNotFound
	}

	// dst must be a pointer
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		return ErrDstMustBePointer
	}

	// kind
	dstKind := kind.KindOf(dst)

	d.lock()
	defer d.unlock()

	for _, r := range d.Data {
		if r.GetID() == i && kind.KindOf(r) == dstKind {
			dstValue.Elem().Set(reflect.ValueOf(r))
			return nil
		}
	}

	return ErrRecordNotFound
}

func (d *DB) Remove(id id.ID, kind kind.Kind) error {
	d.lock()
	defer d.unlock()

	i := d.indexOf(FilterByID(id), FilterByKind(kind))
	if i == -1 {
		return ErrRecordNotFound
	}

	d.Data = append(d.Data[:i], d.Data[i+1:]...)

	return nil
}

func (d *DB) Count() int {
	d.lock()
	defer d.unlock()

	return len(d.Data)
}

func (d *DB) CountByKind(k kind.Kind) int {
	d.lock()
	defer d.unlock()

	var count int

	for _, r := range d.Data {
		if kind.KindOf(r) == k {
			count++
		}
	}

	return count
}
