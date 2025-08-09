package word_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWord(t *testing.T) {
	w := word.New(word.Adjective, "azul", "cielo")
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

	w.ID = db.ID(5)
	assert.Equal(t, db.ID(5), w.GetID())

	s := w.SetID(db.ID(10))
	assert.Equal(t, db.ID(10), s.GetID())

	assert.Equal(t, kind.Word, kind.KindOf(w))
	assert.Equal(t, word.Adjective, w.Type)
}
