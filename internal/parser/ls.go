package parser

import (
	"slices"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
)

type Ordinal int

const (
	First  Ordinal = 1
	Second Ordinal = 2
)

type LS struct {
	original string
	words    []*word.Word
	sub      bool
}

func NewLS() LS {
	return LS{words: []*word.Word{}}
}

func (ls LS) Original() string {
	return ls.original
}

func (ls LS) IsSub() bool {
	return ls.sub
}

// IsEmpty checks if the LS has no words.
func (ls LS) isEmpty() bool {
	return len(ls.words) == 0
}

func (ls *LS) addWord(w *word.Word) {
	if w == nil {
		return
	}

	validTypes := []word.WordType{word.Adverb, word.Verb, word.Noun, word.Adjective}
	if !slices.Contains(validTypes, w.Type) {
		return
	}

	ls.words = append(ls.words, w)
}

func (ls LS) Has(t word.WordType) bool {
	return ls.Get(t, First) != nil
}

func (ls *LS) setVerb(w *word.Word) {
	if ls.Has(word.Verb) {
		return
	}

	words := []*word.Word{}

	for _, ww := range ls.words {
		if ww.Type == word.Adjective || ww.Type == word.Noun {
			words = append(words, w)
			words = append(words, ww)

			continue
		}

		words = append(words, ww)
	}

	ls.words = words
}

func (ls LS) Get(t word.WordType, ord Ordinal) *word.Word {
	count := 0
	for _, w := range ls.words {
		if w.Type == t {
			count++

			if count == int(ord) {
				return w
			}
		}
	}

	return nil
}

func (ls LS) GetIndexOf(t word.WordType, ord Ordinal) int {
	count := 0
	for i, w := range ls.words {
		if w.Type == t {
			count++
			if count == int(ord) {
				return i
			}
		}
	}

	return -1
}

func (ls LS) String() string {
	tokens := []string{}
	for _, w := range ls.words {
		tokens = append(tokens, w.Synonyms[0])
	}

	return strings.Join(tokens, " ")
}
