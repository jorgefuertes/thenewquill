package db

import (
	"reflect"
	"sync"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/log"
)

const NotFound int = -1

type Allow bool

const (
	AllowNoID     Allow = true
	DontAllowNoID Allow = false
)

type Storeable interface {
	GetID() ID
	SetID(id ID) Storeable
	Validate(allowNoID Allow) error
}

type DB struct {
	nextID ID
	mut    *sync.Mutex
	Labels []Label
	Data   []Storeable
}

func New() *DB {
	return &DB{
		mut:    &sync.Mutex{},
		nextID: MinMeaningfulID,
		Labels: []Label{UndefinedLabel, ConfigLabel, UnderscoreLabel, WildcardLabel},
		Data:   make([]Storeable, 0),
	}
}

func (d *DB) Len() int {
	d.mut.Lock()
	defer d.mut.Unlock()

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
	d.mut.Lock()
	defer d.mut.Unlock()

	return d.indexOf(filters...) != NotFound
}

func (d *DB) Create(labelName string, s Storeable) (ID, error) {
	if s.GetID() != UndefinedLabel.ID {
		return UndefinedLabel.ID, ErrCannotCreateWithDefinedID
	}

	if err := s.Validate(AllowNoID); err != nil {
		return UndefinedLabel.ID, err
	}

	label, err := d.AddLabel(labelName, kind.KindOf(s) == kind.Variable)
	if err != nil {
		return UndefinedLabel.ID, err
	}

	s = s.SetID(label.ID)

	return label.ID, d.Append(s)
}

func (d *DB) Append(s Storeable) error {
	if kind.KindOf(s) == kind.None {
		return ErrKindCannotBeNone
	}

	if d.Exists(FilterByID(s.GetID()), FilterByKind(kind.KindOf(s))) {
		l, err := d.GetLabel(s.GetID())
		if err != nil {
			l = Label{Name: err.Error()}
		}

		k := kind.KindOf(s)
		log.Error("❗ duplicated record, %s %q", k, l.Name)

		old := s
		_ = d.Get(s.GetID(), &old)
		log.Error("❗ previous record kind %q: %v", kind.KindOf(old), old)

		return ErrDuplicatedRecord
	}

	if err := s.GetID().Validate(true); err != nil {
		return err
	}

	if err := s.Validate(AllowNoID); err != nil {
		return err
	}

	d.mut.Lock()
	defer d.mut.Unlock()

	d.Data = append(d.Data, s)

	return nil
}

func (d *DB) Update(s Storeable) error {
	if !d.Exists(FilterByID(s.GetID()), FilterByKind(kind.KindOf(s))) {
		return ErrRecordNotFound
	}

	d.mut.Lock()
	defer d.mut.Unlock()

	idx := d.indexOf(FilterByID(s.GetID()), FilterByKind(kind.KindOf(s)))

	d.Data = append(d.Data[:idx], d.Data[idx+1:]...)
	d.Data = append(d.Data, s)

	return nil
}

func (d *DB) GetByLabel(labelName string, dst any) error {
	label, err := d.GetLabelByName(labelName)
	if err != nil {
		return err
	}

	return d.Get(label.ID, dst)
}

func (d *DB) Get(id ID, dst any) error {
	if id < MinMeaningfulID {
		return ErrRecordNotFound
	}

	// dst must be a pointer
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		return ErrDstMustBePointer
	}

	// kind
	dstKind := kind.KindOf(dst)

	d.mut.Lock()
	defer d.mut.Unlock()

	for _, r := range d.Data {
		if r.GetID() == id && kind.KindOf(r) == dstKind {
			dstValue.Elem().Set(reflect.ValueOf(r))
			return nil
		}
	}

	return ErrRecordNotFound
}

func (d *DB) Remove(id ID, kind kind.Kind) error {
	d.mut.Lock()
	defer d.mut.Unlock()

	i := d.indexOf(FilterByID(id), FilterByKind(kind))
	if i == -1 {
		return ErrRecordNotFound
	}

	d.Data = append(d.Data[:i], d.Data[i+1:]...)

	return nil
}

func (d *DB) Reset() {
	d.mut.Lock()
	defer d.mut.Unlock()

	// reset labels
	d.Labels = []Label{UndefinedLabel, ConfigLabel, UnderscoreLabel, WildcardLabel}
	d.nextID = MinMeaningfulID

	// reset entities
	d.Data = make([]Storeable, 0)
}

func (d *DB) Count() int {
	d.mut.Lock()
	defer d.mut.Unlock()

	return len(d.Data)
}

func (d *DB) CountByKind(k kind.Kind) int {
	d.mut.Lock()
	defer d.mut.Unlock()

	var count int

	for _, r := range d.Data {
		if kind.KindOf(r) == k {
			count++
		}
	}

	return count
}
