package kind_test

import (
	"fmt"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"

	"github.com/stretchr/testify/assert"
)

func TestHumanName(t *testing.T) {
	testCases := []struct {
		k    kind.Kind
		want string
	}{
		{kind.None, "None"},
		{kind.Label, "Labels"},
		{kind.Param, "Config"},
		{kind.Variable, "Variables"},
		{kind.Word, "Vocabulary"},
		{kind.Message, "Messages"},
		{kind.Item, "Items"},
		{kind.Location, "Locations"},
		{kind.Character, "Characters"},
		{kind.Table, "Tables"},
		{kind.Process, "Processes"},
		{kind.Blob, "Binary Objects"},
		{kind.Test, "Test Objects"},
	}

	for _, tc := range testCases {
		t.Run(tc.want, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.k.HumanName())
		})
	}
}

func TestHumanNameUnknownFallsBackToNone(t *testing.T) {
	assert.Equal(t, "None", kind.Kind(99).HumanName())
}

func TestInt(t *testing.T) {
	assert.Equal(t, uint8(0), kind.None.Int())
	assert.Equal(t, uint8(6), kind.Item.Int())
	assert.Equal(t, kind.Kind(uint8(7)), kind.Location)
}

func TestIs(t *testing.T) {
	t.Run("matches canonical name", func(t *testing.T) {
		assert.True(t, kind.Item.Is("item"))
		assert.True(t, kind.Location.Is("location"))
	})

	t.Run("matches the integer value as string", func(t *testing.T) {
		assert.True(t, kind.Item.Is(fmt.Sprint(kind.Item.Int())))
	})

	t.Run("rejects mismatches", func(t *testing.T) {
		assert.False(t, kind.Item.Is("character"))
		assert.False(t, kind.Item.Is("99"))
	})
}
