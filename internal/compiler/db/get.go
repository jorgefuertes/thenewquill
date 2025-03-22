package db

import "thenewquill/internal/compiler/section"

func (db *DB) GetRecord(sec section.Section, label string) *Record {
	for _, r := range db.Records {
		if r.Section == sec && r.Label == label {
			return &r
		}
	}

	return nil
}
