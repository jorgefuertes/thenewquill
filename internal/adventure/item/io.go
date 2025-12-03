package item

import (
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/util"
)

const (
	KindFieldIndex = iota
	IDFieldIndex
	NounFieldIndex
	AdjFieldIndex
	DescFieldIndex
	WeightFieldIndex
	MaxWeightFieldIndex
	ContainerFieldIndex
	WearableFieldIndex
	CreatedFieldIndex
	WornFieldIndex
)

func (i Item) Export() string {
	return fmt.Sprintf("%d|%d|%d|%d|%s|%d|%d|%s|%s|%s|%s",
		kind.Item,
		i.ID,
		i.NounID,
		i.AdjectiveID,
		util.EscapeField(i.Description),
		i.Weight,
		i.MaxWeight,
		util.BoolToString(i.Container),
		util.BoolToString(i.Wearable),
		util.BoolToString(i.Created),
		util.BoolToString(i.Worn),
	)
}

func Import(s string) (Item, error) {
	fields := util.SplitIntoFields(s)

	i := Item{}
	i.ID = primitive.IDFromString(fields[IDFieldIndex])
	i.NounID = primitive.IDFromString(fields[NounFieldIndex])
	i.AdjectiveID = primitive.IDFromString(fields[AdjFieldIndex])
	i.Description = fields[DescFieldIndex]
	i.Weight = util.StringToInt(fields[WeightFieldIndex])
	i.MaxWeight = util.StringToInt(fields[MaxWeightFieldIndex])
	i.Container = util.StringToBool(fields[ContainerFieldIndex])
	i.Wearable = util.StringToBool(fields[WearableFieldIndex])
	i.Created = util.StringToBool(fields[CreatedFieldIndex])
	i.Worn = util.StringToBool(fields[WornFieldIndex])

	return i, nil
}
