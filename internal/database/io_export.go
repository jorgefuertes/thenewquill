package database

import (
	"encoding/base64"
	"fmt"
	"time"
)

const (
	version      = "1.0"
	databaseType = "db"
	saveType     = "save"
	labelBegin   = "L:"
	recordBegin  = "R:"
)

func (db *DB) Export(filename string) (int, int, error) {
	params := db.getParams()

	f, err := createFile(filename,
		"The New Quill Adventure Database",
		fmt.Sprintf("Format version: %s, type: %s", version, databaseType),
		fmt.Sprintf("Labels: %d", db.CountLabels()),
		fmt.Sprintf("Records: %d", db.CountRecords()),
		fmt.Sprintf("Timestamp: %d", time.Now().Unix()),
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

	// labels
	for id, label := range db.labels {
		if id < minLabelID {
			continue
		}

		if err := f.writeLn("%s%d|%s", labelBegin, id, label); err != nil {
			return 0, 0, err
		}
	}

	// records
	for id, r := range db.data {
		data := base64.StdEncoding.EncodeToString(r.Data)
		if err := f.writeLn("%s%d|%d|%d|%s", recordBegin, id, r.LabelID, r.Kind, data); err != nil {
			return 0, 0, err
		}
	}

	return f.close()
}
