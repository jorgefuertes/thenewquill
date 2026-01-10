package database

import (
	"encoding/binary"
	"fmt"
	"time"
)

const (
	databaseType  byte = 'D'
	exportVersion byte = 1
)

var endian = binary.LittleEndian

// Export exports the database to a file
// returns the number of bytes sent (before compression) and the final number of bytes in the file
func (db *DB) Export(filename string) (int64, int64, error) {
	params := db.getParams()

	// human header
	f, err := newFileWriter(filename, databaseType, exportVersion,
		"The New Quill Adventure Database",
		fmt.Sprintf("Format version: %d, type: %c", exportVersion, databaseType),
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

	// labels count
	f.write(db.CountLabels())

	// labels
	for id, label := range db.labels {
		if len(label) > 255 {
			return 0, 0, fmt.Errorf("label %d is too long (%d bytes)", id, len(label))
		}

		f.write(id)
		f.write(uint8(len(label))) // label length
		f.write([]byte(label))
	}

	// records count
	f.write(uint32(len(db.data)))

	// records
	for id, record := range db.data {
		f.write(id)
		f.write(record.LabelID)
		f.write(record.Kind.Byte())
		f.write(uint64(len(record.Data))) // data length
		f.write(record.Data)
	}

	return f.close()
}
