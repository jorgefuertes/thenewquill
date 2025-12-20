package word

import (
	"slices"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/database/primitive"
	"github.com/jorgefuertes/thenewquill/internal/util"
)

type Word struct {
	ID       uint32
	LabelID  uint32
	Type     WordType
	Synonyms []string
}

var _ adapter.Storeable = &Word{}

func New(labelID uint32, t WordType, synonyms ...string) *Word {
	for i, s := range synonyms {
		synonyms[i] = strings.ToLower(s)
	}

	return &Word{ID: primitive.UndefinedID, LabelID: labelID, Type: t, Synonyms: synonyms}
}

func (w *Word) SetID(id uint32) {
	w.ID = id
}

func (w Word) GetID() uint32 {
	return w.ID
}

func (w *Word) SetLabelID(id uint32) {
	w.ID = id
}

func (w Word) GetLabelID() uint32 {
	return w.ID
}

func (w Word) HasSynonym(syn string) bool {
	syn = strings.ToLower(syn)

	// check for exact match
	if slices.Contains(w.Synonyms, syn) {
		return true
	}

	// check without accent or symbols
	syn = util.RemoveAccents(syn)
	syn = util.RemoveSymbols(syn)

	return slices.Contains(w.Synonyms, syn)
}

func (w Word) Is(t WordType, syn string) bool {
	if w.Type != t {
		return false
	}

	syn = strings.ToLower(syn)
	if slices.Contains(w.Synonyms, syn) {
		return true
	}

	// check without accent or symbols
	syn = util.RemoveAccents(syn)
	syn = util.RemoveSymbols(syn)

	return slices.Contains(w.Synonyms, syn)
}
