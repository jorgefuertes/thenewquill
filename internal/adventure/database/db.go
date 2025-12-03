package database

import (
	"math"
	"reflect"
	"sync"

	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

const NotFound int = -1

type DB struct {
	nextID primitive.ID
	carry  bool
	mut    *sync.Mutex
	labels []primitive.Label
	data   []adapter.Storeable
}

func (db *DB) lock() {
	db.mut.Lock()
}

func (db *DB) unlock() {
	db.mut.Unlock()
}

func New() *DB {
	return &DB{
		mut:    &sync.Mutex{},
		nextID: primitive.MinID,
		labels: []primitive.Label{primitive.Underscore, primitive.Wildcard},
		data:   []adapter.Storeable{},
	}
}

func (db *DB) GetNextID() (primitive.ID, error) {
	db.lock()
	defer db.unlock()

	if db.nextID == math.MaxUint32 && db.carry {
		return primitive.UndefinedID, ErrLimitReached
	}

	id := db.nextID

	if db.nextID == math.MaxUint32 {
		db.carry = true
	} else {
		db.nextID++
	}

	return id, nil
}

func (db *DB) Reset() {
	db.lock()
	defer db.unlock()

	// reset labels
	db.nextID = primitive.MinID

	// reset entities
	db.data = []adapter.Storeable{}
}

func (db *DB) Len() int {
	db.lock()
	defer db.unlock()

	return len(db.data)
}

func (db *DB) indexOf(filters ...filter) int {
	for i, r := range db.data {
		if db.matches(r, filters...) {
			return i
		}
	}

	return NotFound
}

func (db *DB) Exists(filters ...filter) bool {
	db.lock()
	defer db.unlock()

	return db.indexOf(filters...) != NotFound
}

func (db *DB) Create(s adapter.Storeable) (primitive.ID, error) {
	if s == nil {
		return primitive.UndefinedID, ErrNilStoreable
	}

	if err := s.Validate(true); err != nil {
		return primitive.UndefinedID, err
	}

	if s.GetID() == primitive.UndefinedID {
		i, err := db.GetNextID()
		if err != nil {
			return primitive.UndefinedID, err
		}
		s.SetID(i)
	}

	if db.Exists(FilterByID(s.GetID()), FilterByKind(kind.KindOf(s))) {
		return s.GetID(), ErrDuplicatedRecord
	}

	db.lock()
	defer db.unlock()

	db.data = append(db.data, s)

	return s.GetID(), nil
}

func (db *DB) Update(s adapter.Storeable) error {
	if s == nil {
		return ErrNilStoreable
	}

	if s.GetID() == primitive.UndefinedID {
		return primitive.ErrUndefinedID
	}

	if !db.Exists(FilterByID(s.GetID()), FilterByKind(kind.KindOf(s))) {
		return ErrRecordNotFound
	}

	db.lock()
	defer db.unlock()

	idx := db.indexOf(FilterByID(s.GetID()), FilterByKind(kind.KindOf(s)))

	db.data = append(db.data[:idx], db.data[idx+1:]...)
	db.data = append(db.data, s)

	return nil
}

func (db *DB) Get(i primitive.ID, dst any) error {
	if i < primitive.MinID {
		return ErrRecordNotFound
	}

	// dst must be a pointer
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		return ErrDstMustBePointer
	}

	// kind
	dstKind := kind.KindOf(dst)

	db.lock()
	defer db.unlock()

	for _, r := range db.data {
		if r.GetID() == i && kind.KindOf(r) == dstKind {
			dstValue.Elem().Set(reflect.ValueOf(r).Elem())
			return nil
		}
	}

	return ErrRecordNotFound
}

func (db *DB) Remove(id primitive.ID, kind kind.Kind) error {
	db.lock()
	defer db.unlock()

	i := db.indexOf(FilterByID(id), FilterByKind(kind))
	if i == -1 {
		return ErrRecordNotFound
	}

	db.data = append(db.data[:i], db.data[i+1:]...)

	return nil
}

func (db *DB) Count(filters ...filter) int {
	if len(filters) == 0 {
		db.lock()
		defer db.unlock()

		return len(db.data)
	}

	res := db.Query(filters...)
	defer res.Close()

	return res.Count()
}
