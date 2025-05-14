package character

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

type Character struct {
	ID          db.ID
	NameID      db.ID
	AdjectiveID db.ID
	Description string
	LocationID  db.ID
	Created     bool
	Human       bool
}

var _ db.Storeable = Character{}

func (c Character) GetID() db.ID {
	return c.ID
}

func (c Character) GetKind() (db.Kind, db.SubKind) {
	return db.Chars, db.NoSubKind
}

func New(id db.ID, nameID db.ID, adjectiveID db.ID) Character {
	return Character{
		ID:          id,
		NameID:      nameID,
		AdjectiveID: adjectiveID,
		Description: "",
		LocationID:  db.UndefinedLabel.ID,
		Created:     false,
		Human:       false,
	}
}

func (c Character) Validate() error {
	if c.ID == db.UndefinedLabel.ID {
		return ErrEmptyLabel
	}

	if c.ID < db.MinMeaningfulID {
		return ErrWrongLabel
	}

	if c.Description == "" {
		return ErrEmptyDescription
	}

	return nil
}
