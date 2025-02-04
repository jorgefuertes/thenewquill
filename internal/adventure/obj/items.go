package obj

import (
	"thenewquill/internal/adventure/voc"
	"thenewquill/internal/log"
)

type Items []*Item

func NewItems() Items {
	return Items{}
}

func (i *Items) Len() int {
	return len(*i)
}

func (i *Items) Get(label string) *Item {
	for _, item := range *i {
		if item.label == label {
			return item
		}
	}

	return nil
}

func (i *Items) Exists(label string) bool {
	return i.Get(label) != nil
}

func (i *Items) ExistsNounAdj(noun, adjective *voc.Word) bool {
	for _, item := range *i {
		if item.noun.IsEqual(noun) && item.adjective.IsEqual(adjective) {
			return true
		}
	}

	return false
}

func (i *Items) Add(newItem *Item) error {
	if i.Exists(newItem.label) {
		return ErrDuplicateLabel
	}

	if newItem.label == "" {
		return ErrEmptyLabel
	}

	if newItem.noun == nil {
		return ErrNounCannotBeNil
	}

	if i.ExistsNounAdj(newItem.noun, newItem.adjective) {
		log.Warning(
			"Duplicate noun/adjective for object '%s': %s %s",
			newItem.label,
			newItem.noun.Label,
			newItem.adjective.Label,
		)
	}

	*i = append(*i, newItem)

	return nil
}
