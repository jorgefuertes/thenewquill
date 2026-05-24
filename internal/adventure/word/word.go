package word

import (
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database/adapter"
	"github.com/jorgefuertes/thenewquill/internal/util"
)

const MaxSynonymLen = 5

type Word struct {
	ID           uint32
	LabelID      uint32   `valid:"required"`
	Type         WordType `valid:"required"`
	Synonyms     []string `valid:"count(1|50),len(1|25)"`
	IsConnection bool
	IsItem       bool
	IsCharacter  bool
}

var _ adapter.Storeable = &Word{}

func New(labelID uint32, t WordType, synonyms ...string) *Word {
	seen := make(map[string]struct{}, len(synonyms))
	deduped := make([]string, 0, len(synonyms))

	for _, s := range synonyms {
		s = strings.ToLower(s)
		if t == Verb {
			s = util.TruncateRunes(s, MaxSynonymLen)
		}

		if _, ok := seen[s]; ok {
			continue
		}

		seen[s] = struct{}{}
		deduped = append(deduped, s)
	}

	return &Word{ID: 0, LabelID: labelID, Type: t, Synonyms: deduped}
}

func (w Word) GetKind() kind.Kind {
	return kind.Word
}

func (w *Word) SetID(id uint32) {
	w.ID = id
}

func (w Word) GetID() uint32 {
	return w.ID
}

func (w *Word) SetLabelID(id uint32) {
	w.LabelID = id
}

func (w Word) GetLabelID() uint32 {
	return w.LabelID
}

func (w Word) HasSynonym(syn string) bool {
	if w.Type == Verb {
		syn = util.TruncateRunes(syn, MaxSynonymLen)
	}

	return util.ContainsString(w.Synonyms, syn)
}

func (w Word) Is(t WordType, syn string) bool {
	if w.Type != t {
		return false
	}

	if t == Verb {
		syn = util.TruncateRunes(syn, MaxSynonymLen)
	}

	return util.ContainsString(w.Synonyms, syn)
}
