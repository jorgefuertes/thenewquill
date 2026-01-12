package parser_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/internal/parser"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	db := database.NewDB()
	wordStore := word.NewService(db)

	conjs := []string{"and", "then", "also"}
	for _, c := range conjs {
		labelID, err := db.CreateLabel(c)
		require.NoError(t, err)
		require.NotZero(t, labelID)

		w := word.New(labelID, word.Conjunction, c)
		id, err := wordStore.Create(w)
		require.NoError(t, err)
		require.NotZero(t, id)
	}

	p, err := parser.New(wordStore, parser.EN)
	require.NoError(t, err)
	require.NotNil(t, p)

	t.Run("Parse", func(t *testing.T) {
		input := "take the key and open the door. then go north!"
		p.Parse(input)
	})
}
