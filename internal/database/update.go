package database

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

func (db *DB) Update(entity any) error {
	id, labelID, err := checkEntity(entity)
	if err != nil {
		return err
	}

	if id == 0 {
		return ErrMissingIDToUpdate
	}

	if !db.Exists(id) {
		return ErrRecordNotFound
	}

	if !db.ExistsLabelID(labelID) {
		return ErrLabelNotFound
	}

	_, r, ok := db.Query(FilterByID(id), FilterByKind(kind.KindOf(entity))).getFirstRecord()
	if !ok {
		return ErrRecordNotFound
	}

	if r.Kind != kind.KindOf(entity) {
		return ErrWrongUpdateKind
	}

	r.LabelID = labelID

	r.Data, err = cbor.Marshal(entity)
	if err != nil {
		return err
	}

	if db.IsFrozen() {
		// add to last snapshot
		db.snapshots[len(db.snapshots)-1][id] = r

		return nil
	}

	// add to main data
	db.data[id] = r

	return nil
}
