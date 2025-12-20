package database

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

func (db *DB) Create(entity any) (uint32, error) {
	if db.IsFrozen() {
		return 0, ErrDatabaseIsFrozen
	}

	if db.isFullOfRecords() {
		return 0, ErrDatabaseIsFull
	}

	id, labelID, err := checkEntity(entity)
	if err != nil {
		return id, err
	}

	if id != 0 {
		return id, ErrIDFieldIsNotZero
	}

	if !db.ExistsLabelID(labelID) {
		return id, ErrLabelNotFound
	}

	db.lock()
	defer db.unlock()

	r := Record{
		LabelID: labelID,
		Kind:    kind.KindOf(entity),
		Data:    []byte{},
	}

	id, err = db.getNewDataID()
	if err != nil {
		return 0, err
	}

	setID(entity, id)

	r.Data, err = cbor.Marshal(entity)
	if err != nil {
		return 0, err
	}

	db.data[id] = r

	return id, nil
}
