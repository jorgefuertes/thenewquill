package database

import (
	"encoding/binary"

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

	if f.version != saveVersion {
		log.Fatal("Invalid database header %d, expected version %d", f.version, saveVersion)
	}

	if f.fileType != saveType {
		return ErrInvalidFormatHeader
	}

	db.lock()
	defer db.unlock()

	// snapshots count
	var snapsCount byte
	if err := binary.Read(f.z, endian, &snapsCount); err != nil {
		return err
	}

	db.snapshots = make(snapshots, 0, snapsCount)

	for i := byte(0); i < snapsCount; i++ {
		var recordsCount uint32
		if err := binary.Read(f.z, endian, &recordsCount); err != nil {
			return err
		}

		snap := make(map[uint32]Record, recordsCount)

		for j := uint32(0); j < recordsCount; j++ {
			var id uint32
			if err := binary.Read(f.z, endian, &id); err != nil {
				return err
			}

			var r Record
			if err := binary.Read(f.z, endian, &r.LabelID); err != nil {
				return err
			}
			if err := binary.Read(f.z, endian, &r.Kind); err != nil {
				return err
			}

			var dataLen uint64
			if err := binary.Read(f.z, endian, &dataLen); err != nil {
				return err
			}

			r.Data = make([]byte, dataLen)
			if err := binary.Read(f.z, endian, &r.Data); err != nil {
				return err
			}

			snap[id] = r
		}

		db.snapshots = append(db.snapshots, snap)
	}

	return nil
}
