package database

import (
	"slices"

	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
)

func (db *DB) GetLabelIDForStoreable(id primitive.ID) (primitive.ID, error) {
	db.lock()
	defer db.unlock()

	for _, r := range db.data {
		if r.GetID() == id {
			return r.GetLabelID(), nil
		}
	}

	return primitive.UndefinedID, ErrNotFound
}

func (db *DB) GetLabelForStoreable(id primitive.ID) (primitive.Label, error) {
	labelID, err := db.GetLabelIDForStoreable(id)
	if err != nil {
		return primitive.Label(""), err
	}

	return db.GetLabel(labelID)
}

func (db *DB) GetByLabel(labelOrString any, dst any) error {
	label, err := primitive.LabelFromLabelOrString(labelOrString)
	if err != nil {
		return err
	}

	labelID, err := db.GetLabelID(label)
	if err != nil {
		return err
	}

	return db.Get(labelID, dst)
}

func (db *DB) GetLabel(id primitive.ID) (primitive.Label, error) {
	db.lock()
	defer db.unlock()

	if int(id) >= len(db.labels) {
		return primitive.Label(""), ErrNotFound
	}

	return db.labels[id], nil
}

func (db *DB) GetLabelOrBlank(id primitive.ID) primitive.Label {
	label, _ := db.GetLabel(id)

	return label
}

func (db *DB) GetLabelID(labelOrString any) (primitive.ID, error) {
	label, err := primitive.LabelFromLabelOrString(labelOrString)
	if err != nil {
		return primitive.UndefinedID, err
	}

	db.lock()
	defer db.unlock()

	i, ok := slices.BinarySearch(db.labels, label)
	if ok {
		return primitive.ID(i), nil
	}

	return primitive.UndefinedID, ErrNotFound
}

func (db *DB) CreateLabelIfNotExists(labelOrString any, allowComposite bool) (primitive.ID, error) {
	label, err := primitive.LabelFromLabelOrString(labelOrString)
	if err != nil {
		return primitive.UndefinedID, err
	}

	if err := label.Validate(allowComposite); err != nil {
		return primitive.UndefinedID, err
	}

	i, ok := slices.BinarySearch(db.labels, label)
	if ok {
		return primitive.ID(i), nil
	}

	db.lock()
	defer db.unlock()

	id := primitive.ID(len(db.labels))
	db.labels = append(db.labels, label)

	return id, nil
}
