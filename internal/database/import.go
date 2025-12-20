package database

import (
	"encoding/base64"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/pkg/log"
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

	labelRg := regexp.MustCompile(`^L:(\d+)\|([\d\p{L}\-_\.]+)$`)
	recordRg := regexp.MustCompile(`^R:(\d+)\|(\d+)\|(\d+)\|(.*)$`)

	for {
		line, err := f.readLn()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if labelRg.MatchString(line) {
			matches := labelRg.FindAllStringSubmatch(line, 1)

			if len(matches[0]) < 3 {
				return fmt.Errorf("invalid label line: %s", line)
			}

			i, err := strconv.Atoi(matches[0][1])
			if err != nil {
				return err
			}

			id := uint32(i)
			db.labels[id] = matches[0][2]

			if id > db.lastLabelID {
				db.lastLabelID = id
			}

			continue
		}

		if recordRg.MatchString(line) {
			matches := recordRg.FindAllStringSubmatch(line, 1)

			if len(matches[0]) < 5 {
				return fmt.Errorf("invalid record line: %s", line)
			}

			i, err := strconv.Atoi(matches[0][1])
			if err != nil {
				return err
			}

			id := uint32(i)
			if id > db.lastDataID {
				db.lastDataID = id
			}

			labelID, err := strconv.Atoi(matches[0][2])
			if err != nil {
				return err
			}

			k, err := strconv.Atoi(matches[0][3])
			if err != nil {
				return err
			}

			data, err := base64.StdEncoding.DecodeString(matches[0][4])
			if err != nil {
				return err
			}

			r := Record{
				LabelID: uint32(labelID),
				Kind:    kind.KindFromByte(byte(k)),
				Data:    []byte(data),
			}

			db.data[id] = r

			continue
		}

		log.Warning("Skipping unknown line: %s", line)
	}

	return nil
}
