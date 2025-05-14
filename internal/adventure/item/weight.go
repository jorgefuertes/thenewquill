package item

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

func TotalWeightOf(i Item, d *db.DB) (int, error) {
	if !i.IsContainer {
		return i.Weight, nil
	}

	w := i.Weight
	for _, containedID := range i.Contains {
		var containedItem Item

		if err := d.GetAs(containedID, &containedItem); err != nil {
			return w, err
		}

		w2, err := TotalWeightOf(containedItem, d)
		if err != nil {
			return w, err
		}

		w += w2
	}

	return w, nil
}
