package word_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWord(t *testing.T) {
	w := word.New(0, word.Adjective, "azul", "cielo")
	require.NotEmpty(t, w)

	assert.True(t, w.HasSynonym("azul"))
	assert.True(t, w.HasSynonym("cielo"))
	assert.False(t, w.HasSynonym("negro"))
	assert.False(t, w.HasSynonym("az"))

	assert.True(t, w.Is(word.Adjective, "azul"))
	assert.True(t, w.Is(word.Adjective, "cielo"))
	assert.True(t, w.Is(word.Adjective, "ciélo"))
	assert.False(t, w.Is(word.Adjective, "ciel"))
	assert.False(t, w.Is(word.Verb, "cielo"))

	w.ID = uint32(5)
	assert.Equal(t, uint32(5), w.GetID())

	w.SetID(uint32(10))
	assert.Equal(t, uint32(10), w.GetID())

	assert.Equal(t, kind.Word, kind.KindOf(w))
	assert.Equal(t, word.Adjective, w.Type)
}

func TestWordVerbSynonymTruncation(t *testing.T) {
	w := word.New(0, word.Verb, "examinar", "examina", "exam", "ex")
	require.NotEmpty(t, w)

	// verbs: "examinar" and "examina" both truncate to "exami", deduplicated
	assert.Equal(t, []string{"exami", "exam", "ex"}, w.Synonyms)

	// searching with long words truncates to 5 chars for verbs
	assert.True(t, w.HasSynonym("examinar"))
	assert.True(t, w.HasSynonym("examinando"))
	assert.True(t, w.HasSynonym("exami"))
	assert.True(t, w.HasSynonym("exam"))
	assert.True(t, w.HasSynonym("ex"))
	assert.False(t, w.HasSynonym("e"))
}

func TestWordNonVerbSynonymsNotTruncated(t *testing.T) {
	n := word.New(0, word.Noun, "cuchillo", "navaja")
	require.NotEmpty(t, n)

	// nouns are NOT truncated
	assert.Equal(t, []string{"cuchillo", "navaja"}, n.Synonyms)
	assert.True(t, n.HasSynonym("cuchillo"))
	assert.True(t, n.HasSynonym("navaja"))
	assert.False(t, n.HasSynonym("cuchi"))

	a := word.New(0, word.Adjective, "plateada", "brillante")
	assert.Equal(t, []string{"plateada", "brillante"}, a.Synonyms)
}
