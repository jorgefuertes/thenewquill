package word

import (
	"slices"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/id"
	"github.com/jorgefuertes/thenewquill/internal/util"
)

type Word struct {
	ID       id.ID
	Type     WordType
	Synonyms []string
}

var _ adapter.Storeable = Word{}

func New(t WordType, synonyms ...string) Word {
	for i, s := range synonyms {
		synonyms[i] = strings.ToLower(s)
	}

	return Word{ID: id.Undefined, Type: t, Synonyms: synonyms}
}

func (w Word) SetID(id id.ID) adapter.Storeable {
	w.ID = id

	return w
}

func (w Word) GetID() id.ID {
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
