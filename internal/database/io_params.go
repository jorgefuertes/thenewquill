package database

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database/adapter"
)

type val struct {
	ID      uint32
	LabelID uint32
	V       string
}

var _ adapter.Storeable = &val{}

func (v val) GetKind() kind.Kind {
	return kind.Param
}

func (v val) GetID() uint32 {
	return v.ID
}

func (v *val) SetID(id uint32) {
	v.ID = id
}

func (v val) GetLabelID() uint32 {
	return v.LabelID
}

func (v *val) SetLabelID(id uint32) {
	v.LabelID = id
}

func (db *DB) getParams() map[string]string {
	db.lock()
	defer db.unlock()

	params := make(map[string]string, 0)

	for _, r := range db.data {
		if r.Kind == kind.Param {
			label, err := db.getLabel(r.LabelID)
			if err != nil {
				continue
			}

			v := val{}
			err = r.Unmarshal(&v)
			if err != nil {
				continue
			}

			params[label] = v.V
		}
	}

	return params
}
