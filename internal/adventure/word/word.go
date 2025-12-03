package word

import (
	"slices"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
	"github.com/jorgefuertes/thenewquill/internal/util"
)

type Word struct {
	ID       primitive.ID
	LabelID  primitive.ID
	Type     WordType
	Synonyms []string
}

var _ adapter.Storeable = &Word{}

func New(t WordType, synonyms ...string) *Word {
	for i, s := range synonyms {
		synonyms[i] = strings.ToLower(s)
	}

	return &Word{ID: primitive.UndefinedID, Type: t, Synonyms: synonyms}
}

func (w *Word) SetID(id primitive.ID) {
	w.ID = id
}

func (w Word) GetID() primitive.ID {
	return w.ID
}

func (w *Word) SetLabelID(id primitive.ID) {
	w.ID = id
}

func (w Word) GetLabelID() primitive.ID {
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
