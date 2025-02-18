package voc

import (
	"slices"
	"strings"
	"thenewquill/internal/util"
)

type Word struct {
	Label    string
	Type     WordType
	Synonyms []string
}

func (w Word) String() string {
	return w.Label
}

func (w Word) Is(labelOrSynonym string) bool {
	labelOrSynonym = strings.ToLower(labelOrSynonym)

	// check for exact match
	if w.Label == labelOrSynonym || slices.Contains(w.Synonyms, labelOrSynonym) {
		return true
	}

	// check without accent or symbols
	labelOrSynonym = util.RemoveAccents(labelOrSynonym)
	labelOrSynonym = util.RemoveSymbols(labelOrSynonym)

	return w.Label == labelOrSynonym || slices.Contains(w.Synonyms, labelOrSynonym)
}

func (w Word) IsEqual(w2 *Word) bool {
	return w.Label == w2.Label && w.Type == w2.Type
}
