package database

import (
	"fmt"
	"time"
)

func (db *DB) Save(filename string) (int, int, error) {
	if len(db.snapshots) == 0 {
		return 0, 0, ErrNothingToSave
	}

	params := db.getParams()

	f, err := createFile(filename,
		"The New Quill Adventure Save",
		fmt.Sprintf("Format version: %s, type: %s", version, saveType),
		fmt.Sprintf("Timestamp: %d", time.Now().Unix()),
		fmt.Sprintf("Snapshots: %d", len(db.snapshots)),
		"",
		fmt.Sprintf("Title: %s", params["title"]),
		fmt.Sprintf("Author: %s", params["author"]),
		fmt.Sprintf("Version: %s Lang: %s", params["version"], params["language"]),
	)
	if err != nil {
		return 0, 0, err
	}

	db.lock()
	defer db.unlock()

	for i, snap := range db.snapshots {
		if err := f.writeLn("S:%d", i); err != nil {
			return 0, 0, err
		}

		for id, r := range snap {
			if err := f.writeLn("%s%d|%d|%d|%s", recordBegin, id, r.LabelID, r.Kind, r.Data); err != nil {
				return 0, 0, err
			}
		}
	}

	return 0, 0, nil
}
