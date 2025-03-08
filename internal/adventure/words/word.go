package words

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

func New(label string, t WordType, synonyms ...string) *Word {
	return &Word{Label: label, Type: t, Synonyms: synonyms}
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

func (w Word) IsExactlyEqual(w2 Word) bool {
	if w.Label != w2.Label || w.Type != w2.Type {
		return false
	}

	if len(w.Synonyms) != len(w2.Synonyms) {
		return false
	}

	for _, s := range w.Synonyms {
		if !slices.Contains(w2.Synonyms, s) {
			return false
		}
	}

	return true
}

func (w Word) IsEqual(w2 *Word) bool {
	return w.Label == w2.Label && w.Type == w2.Type
}

func (w Word) export() map[string]any {
	data := map[string]any{
		"label": w.Label,
		"type":  int(w.Type),
	}

	data["synonyms"] = strings.Join(w.Synonyms, ",")

	return data
}
