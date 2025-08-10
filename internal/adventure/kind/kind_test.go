package kind_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/character"
	"github.com/jorgefuertes/thenewquill/internal/adventure/config"
	"github.com/jorgefuertes/thenewquill/internal/adventure/item"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/adventure/location"
	"github.com/jorgefuertes/thenewquill/internal/adventure/message"
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
		{"Config", kind.Param, "config"},
		{"Variable", kind.Variable, "var"},
		{"Word", kind.Word, "word"},
		{"Message", kind.Message, "message"},
		{"Item", kind.Item, "item"},
		{"Location", kind.Location, "location"},
		{"Process", kind.Process, "process table"},
		{"Label", kind.Label, "label"},
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
		{"Config", kind.Param, 2},
		{"Variable", kind.Variable, 3},
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
		{"Label", 1, kind.Label},
		{"Config", 2, kind.Param},
		{"Variables", 3, kind.Variable},
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
		{"Config", "config", kind.Param},
		{"Config alias", "cfg", kind.Param},
		{"Variable", "var", kind.Variable},
		{"Word", "vocabulary", kind.Word},
		{"Invalid", "invalid", kind.None},
		{"Case insensitive", "CONFIG", kind.Param},
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
		{"Config", config.Param{}, kind.Param},
		{"Config pointer", &config.Param{}, kind.Param},
		{"Variable", variable.Variable{}, kind.Variable},
		{"None", nil, kind.None},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, kind.KindOf(tc.input))
		})
	}
}
