package database

import (
	"io"
	"regexp"

	"github.com/jorgefuertes/thenewquill/pkg/log"
)

var (
	labelRg  = regexp.MustCompile(`^L:(\d+)\|([\d\p{L}\-_\.]+)$`)
	recordRg = regexp.MustCompile(`^R:(\d+)\|(\d+)\|(\d+)\|(.*)$`)
)

func (db *DB) Import(filename string) error {
	db.ResetDB()
	db.lock()
	defer db.unlock()

	f, err := newFileReader(filename)
	if err != nil {
		return err
	}

	defer f.close()

	if f.version != version {
		log.Warning("Invalid database header %q, expected version %q", f.version, version)
		log.Warning("A database for another runtime version can lead to in game problems")
	}

	if f.fileType != databaseType {
		return ErrInvalidFormatHeader
	}

	for {
		line, err := f.readLn()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if labelRg.MatchString(line) {
			id, label, err := labelFromLine(line)
			if err != nil {
				return err
			}

			if id > db.lastLabelID {
				db.lastLabelID = id
			}

			db.labels[id] = label

			continue
		}

		if recordRg.MatchString(line) {
			id, r, err := recordFromLine(line)
			if err != nil {
				return err
			}

			if id > db.lastDataID {
				db.lastDataID = id
			}

			db.data[id] = r

			continue
		}

		log.Warning("Skipping unknown line: %s", line)
	}

	return nil
}
