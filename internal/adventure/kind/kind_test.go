package kind_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/character"
	"github.com/jorgefuertes/thenewquill/internal/adventure/config"
	"github.com/jorgefuertes/thenewquill/internal/adventure/item"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/adventure/location"
	"github.com/jorgefuertes/thenewquill/internal/adventure/message"
	"github.com/jorgefuertes/thenewquill/internal/adventure/process"
	"github.com/jorgefuertes/thenewquill/internal/adventure/variable"
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"

	"github.com/stretchr/testify/assert"
)

func TestKinds(t *testing.T) {
	assert.NotEmpty(t, kind.Kinds())
	assert.Contains(t, kind.Kinds(), kind.Param)
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
		{"Label", kind.Label, "label"},
		{"Param", kind.Param, "config"},
		{"Config", kind.Param, "config"},
		{"Variable", kind.Variable, "var"},
		{"Word", kind.Word, "word"},
		{"Message", kind.Message, "message"},
		{"Item", kind.Item, "item"},
		{"Location", kind.Location, "location"},
		{"Table", kind.Table, "table"},
		{"Process", kind.Process, "process"},
		{"Character", kind.Character, "character"},
		{"Test", kind.Test, "test"},
		{"Invalid", kind.Kind(254), "none"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.kind.String())
		})
	}
}

// None Label Param Variable Word Message Item Location Character Process Test
func TestKindByte(t *testing.T) {
	tests := []struct {
		name     string
		kind     kind.Kind
		expected byte
	}{
		{"None", kind.None, 0},
		{"Label", kind.Label, 1},
		{"Config", kind.Param, 2},
		{"Variable", kind.Variable, 3},
		{"Word", kind.Word, 4},
		{"Message", kind.Message, 5},
		{"Item", kind.Item, 6},
		{"Location", kind.Location, 7},
		{"Character", kind.Character, 8},
		{"Table", kind.Table, 9},
		{"Process", kind.Process, 10},
		{"Blob", kind.Blob, 11},
		{"Test", kind.Test, 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.kind.Byte(), "Kind: %s", tt.kind.String())
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
		{"Label", 1, kind.Label},
		{"Config", 2, kind.Param},
		{"Variables", 3, kind.Variable},
		{"Word", 4, kind.Word},
		{"Message", 5, kind.Message},
		{"Item", 6, kind.Item},
		{"Location", 7, kind.Location},
		{"Character", 8, kind.Character},
		{"Table", 9, kind.Table},
		{"Process", 10, kind.Process},
		{"Blob", 11, kind.Blob},
		{"Test", 12, kind.Test},
		{"Invalid High", 255, kind.None},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(
				t,
				tt.expected,
				kind.KindFromByte(tt.input),
				"Input byte: %d, expected kind: %d",
				tt.input,
				tt.expected,
			)
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
		{"Config", "config", kind.Param},
		{"Config alias", "cfg", kind.Param},
		{"Variable", "var", kind.Variable},
		{"Word", "vocabulary", kind.Word},
		{"Message", "message", kind.Message},
		{"Item", "object", kind.Item},
		{"Location", "room", kind.Location},
		{"Character", "player", kind.Character},
		{"Table", "table", kind.Table},
		{"Process", "proc", kind.Process},
		{"Blob", "picture", kind.Blob},
		{"Test", "testitem", kind.Test},
		{"Invalid", "invalid", kind.None},
		{"Case insensitive", "CONFIG", kind.Param},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, kind.KindFromString(tt.input), "Input string: %s", tt.input)
		})
	}
}

func TestKindOf(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected kind.Kind
	}{
		{"Item", &item.Item{}, kind.Item},
		{"Character", &character.Character{}, kind.Character},
		{"Location", &location.Location{}, kind.Location},
		{"Word", &word.Word{}, kind.Word},
		{"Message", &message.Message{}, kind.Message},
		{"Config", &config.Param{}, kind.Param},
		{"Variable", &variable.Variable{}, kind.Variable},
		{"Table", &process.Table{}, kind.Table},
		{"Process", &process.Process{}, kind.Process},
		{"None", nil, kind.None},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, kind.KindOf(tc.input))
		})
	}
}
