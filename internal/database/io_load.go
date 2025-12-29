package database

import (
	"io"
	"regexp"

	"github.com/jorgefuertes/thenewquill/pkg/log"
)

func (db *DB) Load(filename string) error {
	f, err := newFileReader(filename)
	if err != nil {
		return err
	}

	defer func() {
		if err := f.close(); err != nil {
			log.Error("Error closing %q: %s", filename, err)
		}
	}()

	if f.version != version {
		log.Fatal("Invalid database header %q, expected version %q", f.version, version)
	}

	if f.fileType != saveType {
		return ErrInvalidFormatHeader
	}

	db.ResetDB()
	db.lock()
	defer db.unlock()

	beginSnapshotRg := regexp.MustCompile(`^S:(\d+)$`)

	for {
		line, err := f.readLn()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if beginSnapshotRg.MatchString(line) {
			db.snapshots = append(db.snapshots, make(data))
		}

		if recordRg.MatchString(line) {
			id, r, err := recordFromLine(line)
			if err != nil {
				return err
			}

			db.snapshots[len(db.snapshots)-1][id] = r
		}
	}

	return nil
}
