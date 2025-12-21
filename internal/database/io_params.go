package database

import "github.com/jorgefuertes/thenewquill/internal/adventure/kind"

func (db *DB) getParams() map[string]string {
	type val struct {
		V string
	}

	db.lock()
	defer db.unlock()

	params := make(map[string]string, 0)

	for _, r := range db.data {
		if r.Kind == kind.Param {
			label, err := db.GetLabel(r.LabelID)
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
