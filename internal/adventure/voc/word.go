package voc

import "slices"

type Word struct {
	Label    string
	Type     WordType
	Synonyms []string
}

func (w Word) String() string {
	return w.Label
}

func (w Word) Is(labelOrSynonym string) bool {
	return w.Label == labelOrSynonym || slices.Contains(w.Synonyms, labelOrSynonym)
}

func (w Word) IsEqual(w2 *Word) bool {
	return w.Label == w2.Label && w.Type == w2.Type
}
