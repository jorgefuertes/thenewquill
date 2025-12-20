package character

type Character struct {
	ID          uint32
	LabelID     uint32
	NounID      uint32
	AdjectiveID uint32
	Description string
	LocationID  uint32
	Created     bool
	Human       bool
}

func New(id, labelID, nounID, adjectiveID uint32) *Character {
	return &Character{
		ID:          id,
		LabelID:     labelID,
		NounID:      nounID,
		AdjectiveID: adjectiveID,
		Description: "",
		Created:     false,
		Human:       false,
	}
}

func (c Character) GetID() uint32 {
	return c.ID
}

func (c *Character) SetID(id uint32) {
	c.ID = id
}

func (c Character) GetLabelID() uint32 {
	return c.LabelID
}

func (c *Character) SetLabelID(id uint32) {
	c.LabelID = id
}
