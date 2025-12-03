package character

import (
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/util"
)

var _ adapter.Exportable = &Character{}

const (
	KindFieldIndex = iota
	IDFieldIndex
	NounFieldIndex
	AdjFieldIndex
	DescFieldIndex
	LocFieldIndex
	CreatedFieldIndex
	HumanFieldIndex
)

func (c Character) Export() string {
	return fmt.Sprintf(
		"%d|%d|%d|%d|%s|%d|%s|%s",
		kind.KindOf(c).Int(),
		c.ID,
		c.NounID,
		c.AdjectiveID,
		util.EscapeField(c.Description),
		c.LocationID,
		util.BoolToString(c.Created),
		util.BoolToString(c.Human),
	)
}

func Import(s string) Character {
	fields := util.SplitIntoFields(s)

	c := Character{}

	c.ID = primitive.IDFromString(fields[IDFieldIndex])
	c.NounID = primitive.IDFromString(fields[NounFieldIndex])
	c.AdjectiveID = primitive.IDFromString(fields[AdjFieldIndex])
	c.Description = fields[DescFieldIndex]
	c.LocationID = primitive.IDFromString(fields[LocFieldIndex])
	c.Created = util.StringToBool(fields[CreatedFieldIndex])
	c.Human = util.StringToBool(fields[HumanFieldIndex])

	return c
}
