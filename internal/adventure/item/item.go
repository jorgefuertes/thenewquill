package item

import (
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/util"
)

type Item struct {
	ID          db.ID
	NounID      db.ID
	AdjectiveID db.ID
	Description string
	Weight      int
	MaxWeight   int
	Container   bool
	Wearable    bool
	Created     bool
	At          db.ID
	Worn        bool
}

var _ db.Storeable = Item{}

func (i Item) Export() string {
	return fmt.Sprintf("%d|%d|%d|%d|%s|%d|%d|%d|%d|%d|%d|%d\n",
		i.GetKind().Byte(),
		i.ID,
		i.NounID,
		i.AdjectiveID,
		util.EscapeExportString(i.Description),
		i.Weight,
		i.MaxWeight,
		util.BoolToInt(i.Container),
		util.BoolToInt(i.Wearable),
		util.BoolToInt(i.Created),
		i.At,
		util.BoolToInt(i.Worn),
	)
}

func New(nounID db.ID, adjectiveID db.ID) Item {
	return Item{
		ID:          db.UndefinedLabel.ID,
		NounID:      nounID,
		AdjectiveID: adjectiveID,
		Weight:      0,
		MaxWeight:   100,
	}
}

func (i Item) SetID(id db.ID) db.Storeable {
	i.ID = id

	return i
}

func (i Item) GetID() db.ID {
	return i.ID
}

func (i Item) GetKind() db.Kind {
	return db.Items
}
