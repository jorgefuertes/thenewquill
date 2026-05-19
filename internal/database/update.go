package database

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database/adapter"
)

func (db *DB) Update(entity adapter.Storeable) error {
	id := entity.GetID()
	labelID := entity.GetLabelID()

	if id == 0 {
		return ErrMissingIDToUpdate
	}

	if !db.ExistsLabelID(labelID) {
		return ErrLabelNotFound
	}

	r, ok := db.data[id]
	if !ok {
		return ErrRecordNotFound
	}

	if r.Kind != kind.KindOf(entity) {
		return ErrWrongUpdateKind
	}

	var err error
	r.Data, err = cbor.Marshal(entity)
	if err != nil {
		return err
	}

	r.LabelID = labelID

	if db.IsFrozen() {
		// add to last snapshot; create one implicitly if the caller has not
		// opened a transaction yet, so Update never panics on snapshots[-1].
		if len(db.snapshots) == 0 {
			db.snapshots = append(db.snapshots, make(data))
		}
		db.snapshots[len(db.snapshots)-1][id] = r

		return nil
	}

	// add to main data
	db.data[id] = r

	return nil
}
