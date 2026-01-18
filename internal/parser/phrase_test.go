package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPartToPhrases(t *testing.T) {
	t.Run("normal phrase without quotes", func(t *testing.T) {
		input := "  coger la llave  "
		result := partToPhrases(input)

		require.Len(t, result, 1)
		assert.Equal(t, "coger la llave", result[0].str)
		assert.False(t, result[0].isTalking)
	})

	t.Run("normal phrase with mixed case", func(t *testing.T) {
		input := "  CoGeR La LlaVe  "
		result := partToPhrases(input)

		require.Len(t, result, 1)
		assert.Equal(t, "coger la llave", result[0].str)
		assert.False(t, result[0].isTalking)
	})

	t.Run("phrase with double quotes", func(t *testing.T) {
		input := `Say to Elf "Give me the key"`
		result := partToPhrases(input)

		require.Len(t, result, 2)

		// First phrase: command without quotes
		assert.Equal(t, "say to elf", result[0].str)
		assert.False(t, result[0].isTalking)

		// Second phrase: quoted text
		assert.Equal(t, "give me the key", result[1].str)
		assert.True(t, result[1].isTalking)
	})

	t.Run("phrase with single quotes", func(t *testing.T) {
		input := `Say to Elf 'Give me the key'`
		result := partToPhrases(input)

		require.Len(t, result, 2)

		assert.Equal(t, "say to elf", result[0].str)
		assert.False(t, result[0].isTalking)

		assert.Equal(t, "give me the key", result[1].str)
		assert.True(t, result[1].isTalking)
	})

	t.Run("only quoted text - discarded", func(t *testing.T) {
		// Only quoted text without a preceding phrase is discarded
		// Always need at least one phrase with isTalking=false
		input := `"hello world"`
		result := partToPhrases(input)

		assert.Empty(t, result)
	})

	t.Run("multiple quoted sections - only first is processed", func(t *testing.T) {
		// Only the first quoted section is processed; second quotes are treated as literal text
		input := `say "first" and then "second"`
		result := partToPhrases(input)

		require.Len(t, result, 3)

		assert.Equal(t, "say", result[0].str)
		assert.False(t, result[0].isTalking)

		assert.Equal(t, "first", result[1].str)
		assert.True(t, result[1].isTalking)

		assert.Equal(t, `and then "second"`, result[2].str)
		assert.False(t, result[2].isTalking)
	})

	t.Run("text with leading and trailing spaces", func(t *testing.T) {
		input := `  say to elf  "give me the key"  `
		result := partToPhrases(input)

		require.Len(t, result, 2)

		assert.Equal(t, "say to elf", result[0].str)
		assert.False(t, result[0].isTalking)

		assert.Equal(t, "give me the key", result[1].str)
		assert.True(t, result[1].isTalking)
	})

	t.Run("empty string", func(t *testing.T) {
		input := ""
		result := partToPhrases(input)

		assert.Empty(t, result)
	})

	t.Run("only spaces", func(t *testing.T) {
		input := "   "
		result := partToPhrases(input)

		assert.Empty(t, result)
	})

	t.Run("quoted empty string", func(t *testing.T) {
		input := `say ""`
		result := partToPhrases(input)

		// Empty quotes are discarded; only the non-quoted phrase remains
		require.Len(t, result, 1)
		assert.Equal(t, `say`, result[0].str)
		assert.False(t, result[0].isTalking)
	})

	t.Run("curly quotes", func(t *testing.T) {
		input := `say to elf "give me the key"`
		result := partToPhrases(input)

		require.Len(t, result, 2)

		assert.Equal(t, "say to elf", result[0].str)
		assert.False(t, result[0].isTalking)

		assert.Equal(t, "give me the key", result[1].str)
		assert.True(t, result[1].isTalking)
	})
}
