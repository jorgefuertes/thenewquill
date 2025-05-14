package db_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"

	"github.com/stretchr/testify/assert"
)

func TestKinds(t *testing.T) {
	kinds := db.Kinds()
	assert.NotEmpty(t, kinds)
	assert.Contains(t, kinds, db.Config)
	assert.Contains(t, kinds, db.Variables)
	assert.Contains(t, kinds, db.Words)
}

func TestKindString(t *testing.T) {
	tests := []struct {
		name     string
		kind     db.Kind
		expected string
	}{
		{"None", db.None, "none"},
		{"Config", db.Config, "config"},
		{"Variables", db.Variables, "vars"},
		{"Words", db.Words, "words"},
		{"Messages", db.Messages, "messages"},
		{"Items", db.Items, "items"},
		{"Locations", db.Locations, "locations"},
		{"Processes", db.Processes, "process tables"},
		{"Invalid", db.Kind(255), "none"},
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
		kind     db.Kind
		expected byte
	}{
		{"None", db.None, 0},
		{"Config", db.Config, 1},
		{"Variables", db.Variables, 2},
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
		expected db.Kind
	}{
		{"Zero", 0, db.None},
		{"Config", 1, db.Config},
		{"Variables", 2, db.Variables},
		{"Invalid High", 255, db.None},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, db.FromByte(tt.input))
		})
	}
}

func TestFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected db.Kind
	}{
		{"Empty string", "", db.None},
		{"None", "none", db.None},
		{"Unknown", "unknown", db.None},
		{"Config", "config", db.Config},
		{"Config alias", "cfg", db.Config},
		{"Variables", "vars", db.Variables},
		{"Words", "vocabulary", db.Words},
		{"Invalid", "invalid", db.None},
		{"Case insensitive", "CONFIG", db.Config},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, db.FromString(tt.input))
		})
	}
}
