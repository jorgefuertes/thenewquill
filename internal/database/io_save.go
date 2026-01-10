package database

import (
	"fmt"
	"time"
)

const (
	saveType    byte = 'S'
	saveVersion byte = 1
)

func (db *DB) Save(filename string) (int64, int64, error) {
	if len(db.snapshots) == 0 {
		return 0, 0, ErrNothingToSave
	}

	params := db.getParams()

	f, err := newFileWriter(filename, saveType, saveVersion,
		"The New Quill Adventure Save",
		fmt.Sprintf("Format version: %d, type: %c", saveVersion, saveType),
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

	if len(db.snapshots) > 255 {
		return 0, 0, fmt.Errorf("too many snapshots to save (%d), max is 255", len(db.snapshots))
	}

	// snapshots count
	f.write(byte(len(db.snapshots)))

	for _, snap := range db.snapshots {
		// snapshot's records count
		f.write(uint32(len(snap)))

		// records
		for id, r := range snap {
			f.write(id)
			f.write(r.LabelID)
			f.write(r.Kind)
			f.write(uint64(len(r.Data))) // data length
			f.write(r.Data)
		}
	}

	return f.close()
}
