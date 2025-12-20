package database

import (
	"math"
	"regexp"
	"strings"
)

const (
	LabelUnderscore = "_"
	LabelAsterisk   = "*"
)

// GetLabel returns the label if it exists, or ErrLabelNotFound
func (db *DB) GetLabel(id uint32) (string, error) {
	db.mut.Lock()
	defer db.mut.Unlock()

	return db.getLabel(id)
}

// getLabel returns the label if it exists, or ErrLabelNotFound, without locking
func (db *DB) getLabel(id uint32) (string, error) {
	label, ok := db.labels[id]
	if !ok {
		return "", ErrLabelNotFound
	}

	return label, nil
}

// GetLabelID returns the ID of the label if it exists, or ErrLabelNotFound
func (db *DB) GetLabelID(label string) (uint32, error) {
	db.mut.Lock()
	defer db.mut.Unlock()

	return db.getLabelID(label)
}

// getLabelID returns the ID of the label if it exists, or ErrLabelNotFound, without locking
func (db *DB) getLabelID(label string) (uint32, error) {
	label = strings.ToLower(label)

	for i, l := range db.labels {
		if l == label {
			return i, nil
		}
	}

	return 0, ErrLabelNotFound
}

// CreateLabel creates a new label in the database if it does not exist
// returns the ID of the label, created or already existing
func (db *DB) CreateLabel(label string) (uint32, error) {
	if db.isFullOfLabels() {
		return 0, ErrLabelsAreFull
	}

	label = strings.ToLower(label)

	if !regexp.MustCompile(`^([\d\p{L}\-_\.]+|[_\*]{1})$`).MatchString(label) {
		return 0, ErrInvalidLabel
	}

	id, err := db.GetLabelID(label)
	if err == nil {
		return id, nil
	}

	db.mut.Lock()
	defer db.mut.Unlock()

	id, err = db.getNewLabelID()
	if err != nil {
		return 0, err
	}

	db.labels[id] = label

	return id, nil
}

func (db *DB) isFullOfLabels() bool {
	return len(db.labels) == math.MaxUint32
}

func (db *DB) CountLabels() uint32 {
	return uint32(len(db.labels) - 1)
}

func (db *DB) ExistsLabelID(id uint32) bool {
	return id <= db.CountLabels()
}
