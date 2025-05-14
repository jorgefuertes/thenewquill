package words

import (
	"slices"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/util"
)

type Word struct {
	ID       db.ID
	Type     db.SubKind
	Synonyms []string
}

func New(id db.ID, t db.SubKind, synonyms ...string) *Word {
	for i, s := range synonyms {
		synonyms[i] = strings.ToLower(s)
	}

	return &Word{ID: id, Type: t, Synonyms: synonyms}
}

func (w Word) GetID() db.ID {
	return w.ID
}

func (w Word) GetKind() (db.Kind, db.SubKind) {
	return db.Words, w.Type
}

func (w Word) Is(syn string) bool {
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

func (w Word) IsSynonymAndType(syn string, t db.SubKind) bool {
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
