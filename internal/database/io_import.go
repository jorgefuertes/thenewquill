package database

import (
	"encoding/binary"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/pkg/log"
)

func (db *DB) Import(filename string) error {
	f, err := newFileReader(filename)
	if err != nil {
		return err
	}

	defer func() {
		if err := f.close(); err != nil {
			log.Error("Error closing %q: %s", filename, err)
		}
	}()

	if f.version != exportVersion {
		log.Warning("Invalid database header %d, expected version %d", f.version, exportVersion)
		log.Warning("A database for another runtime version can lead to in game problems")
	}

	if f.fileType != databaseType {
		return ErrInvalidFormatHeader
	}

	// reset the database
	db.ResetDB()
	db.lock()
	defer db.unlock()

	// labels count
	var labelsCount uint32
	if err := binary.Read(f.z, endian, &labelsCount); err != nil {
		return err
	}

	// labels
	for i := uint32(1); i <= labelsCount; i++ {
		var id uint32
		if err := binary.Read(f.z, endian, &id); err != nil {
			return err
		}

		var length uint8
		if err := binary.Read(f.z, endian, &length); err != nil {
			return err
		}

		labelBytes := make([]byte, length)
		if err := binary.Read(f.z, endian, &labelBytes); err != nil {
			return err
		}

		db.labels[id] = string(labelBytes)
	}

	// records count
	var recordsCount uint32
	if err := binary.Read(f.z, endian, &recordsCount); err != nil {
		return err
	}

	// records
	for i := uint32(1); i <= recordsCount; i++ {
		var id uint32
		if err := binary.Read(f.z, endian, &id); err != nil {
			return err
		}

		var labelID uint32
		if err := binary.Read(f.z, endian, &labelID); err != nil {
			return err
		}

		var kindByte byte
		if err := binary.Read(f.z, endian, &kindByte); err != nil {
			return err
		}

		var dataLength uint64
		if err := binary.Read(f.z, endian, &dataLength); err != nil {
			return err
		}

		dataBytes := make([]byte, dataLength)
		if err := binary.Read(f.z, endian, &dataBytes); err != nil {
			return err
		}

		db.data[id] = Record{
			LabelID: labelID,
			Kind:    kind.Kind(kindByte),
			Data:    dataBytes,
		}
	}

	return nil
}
