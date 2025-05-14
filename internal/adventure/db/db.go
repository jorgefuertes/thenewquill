package db

import (
	"reflect"
	"sync"
)

const NotFound int = -1

type Storeable interface {
	GetID() ID
	GetKind() (Kind, SubKind)
	Validate() error
}

type DB struct {
	mut    *sync.Mutex
	nextID ID
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

func (d *DB) indexOf(id ID, k Kind, sk SubKind) int {
	for i, r := range d.Data {
		if !isKind(r, k, sk) {
			continue
		}

		if r.GetID() == id {
			return i
		}
	}

	return NotFound
}

func (d *DB) Exists(id ID, k Kind, sk SubKind) bool {
	d.mut.Lock()
	defer d.mut.Unlock()

	return d.indexOf(id, k, sk) != -1
}

func (d *DB) Append(s Storeable) error {
	k, sk := s.GetKind()

	if k == None {
		return ErrKindCannotBeNone
	}

	if sk == AnySubKind {
		return ErrSubKindMustBeDefined
	}

	if d.Exists(s.GetID(), k, sk) {
		return ErrDuplicatedRecord
	}

	d.mut.Lock()
	defer d.mut.Unlock()

	d.Data = append(d.Data, s)

	return nil
}

func (d *DB) Update(s Storeable) error {
	k, sk := s.GetKind()
	if !d.Exists(s.GetID(), k, sk) {
		return ErrRecordNotFound
	}

	d.mut.Lock()
	defer d.mut.Unlock()

	idx := d.indexOf(s.GetID(), k, sk)

	d.Data = append(d.Data[:idx], d.Data[idx+1:]...)
	d.Data = append(d.Data, s)

	return nil
}

func (d *DB) Get(kind Kind, sub SubKind, id ID) (Storeable, error) {
	d.mut.Lock()
	defer d.mut.Unlock()

	for _, r := range d.Data {
		if r.GetID() == id && isKind(r, kind, sub) {
			return r, nil
		}
	}

	return nil, ErrRecordNotFound
}

func (d *DB) GetAs(id ID, kind Kind, sub SubKind, dst Storeable) error {
	r, err := d.Get(kind, sub, id)
	if err != nil {
		return err
	}

	// dst must be a pointer
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		return ErrDstMustBePointer
	}

	dstValue.Elem().Set(reflect.ValueOf(r))

	return nil
}

func (d *DB) GetByKind(k Kind, sk SubKind) []Storeable {
	d.mut.Lock()
	defer d.mut.Unlock()

	var objects []Storeable

	for _, r := range d.Data {
		if isKind(r, k, sk) {
			objects = append(objects, r)
		}
	}

	return objects
}

// GetKind returns all the objects of the given kind in the given slice
func (d *DB) GetByKindAs(kind Kind, sub SubKind, dst interface{}) error {
	// Verify that dst is a pointer to a slice
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr || dstValue.Elem().Kind() != reflect.Slice {
		return ErrDstMustBePointerSlice
	}

	// Get the slice value and the element type
	sliceValue := dstValue.Elem()
	elementType := sliceValue.Type().Elem()

	// Create a new slice
	newSlice := reflect.MakeSlice(sliceValue.Type(), 0, 0)

	objects := d.GetByKind(kind, sub)
	for _, r := range objects {
		// Convert the object to the element type
		if reflect.TypeOf(r).AssignableTo(elementType) {
			newSlice = reflect.Append(newSlice, reflect.ValueOf(r))
		} else {
			return ErrCannotCastFromStoreable
		}
	}

	// Asign the new slice to the destination
	sliceValue.Set(newSlice)

	return nil
}

func (d *DB) Remove(id ID, kind Kind, sub SubKind) error {
	d.mut.Lock()
	defer d.mut.Unlock()

	i := d.indexOf(id, kind, sub)
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

func (d *DB) Count(kind Kind, sub SubKind) int {
	var count int

	d.mut.Lock()
	defer d.mut.Unlock()

	for _, r := range d.Data {
		if isKind(r, kind, sub) {
			count++
		}
	}

	return count
}

func (d *DB) CountAll() int {
	d.mut.Lock()
	defer d.mut.Unlock()

	return len(d.Data)
}
