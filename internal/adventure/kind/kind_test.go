package kind_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/character"
	"github.com/jorgefuertes/thenewquill/internal/adventure/item"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/adventure/location"
	"github.com/jorgefuertes/thenewquill/internal/adventure/message"
	"github.com/jorgefuertes/thenewquill/internal/adventure/variable"
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
	"honnef.co/go/tools/config"

	"github.com/stretchr/testify/assert"
)

func TestKinds(t *testing.T) {
	assert.NotEmpty(t, kind.Kinds())
	assert.Contains(t, kind.Kinds(), kind.Config)
	assert.Contains(t, kind.Kinds(), kind.Variable)
	assert.Contains(t, kind.Kinds(), kind.Word)
}

func TestKindString(t *testing.T) {
	tests := []struct {
		name     string
		kind     kind.Kind
		expected string
	}{
		{"None", kind.None, "none"},
		{"Config", kind.Config, "config"},
		{"Variables", kind.Variable, "vars"},
		{"Words", kind.Word, "words"},
		{"Messages", kind.Message, "messages"},
		{"Items", kind.Item, "items"},
		{"Locations", kind.Location, "locations"},
		{"Processes", kind.Process, "process tables"},
		{"Labels", kind.Label, "labels"},
		{"Invalid", kind.Kind(254), "none"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.kind.String())
		})
	}
}

func TestKindByte(t *testing.T) {
	tests := []struct {
		name     string
		kind     kind.Kind
		expected byte
	}{
		{"None", kind.None, 0},
		{"Config", kind.Config, 1},
		{"Variables", kind.Variable, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.kind.Byte())
		})
	}
}

func TestFromByte(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected kind.Kind
	}{
		{"Zero", 0, kind.None},
		{"Config", 1, kind.Config},
		{"Variables", 2, kind.Variable},
		{"Invalid High", 255, kind.None},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, kind.KindFromByte(tt.input))
		})
	}
}

func TestFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected kind.Kind
	}{
		{"Empty string", "", kind.None},
		{"None", "none", kind.None},
		{"Unknown", "unknown", kind.None},
		{"Config", "config", kind.Config},
		{"Config alias", "cfg", kind.Config},
		{"Variables", "vars", kind.Variable},
		{"Words", "vocabulary", kind.Word},
		{"Invalid", "invalid", kind.None},
		{"Case insensitive", "CONFIG", kind.Config},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, kind.KindFromString(tt.input))
		})
	}
}

func TestKindOf(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected kind.Kind
	}{
		{"Pointer to Item", &item.Item{}, kind.Item},
		{"Item", item.Item{}, kind.Item},
		{"Character", character.Character{}, kind.Character},
		{"Location", location.Location{}, kind.Location},
		{"Word", word.Word{}, kind.Word},
		{"Message", message.Message{}, kind.Message},
		{"Config", config.Config{}, kind.Config},
		{"Variable", variable.Variable{}, kind.Variable},
		{"None", nil, kind.None},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, kind.KindOf(tc.input))
		})
	}
}
