package util_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestValueToString(t *testing.T) {
	tests := []struct {
		input    any
		expected string
	}{
		{input: "example", expected: "example"},
		{input: true, expected: util.TrueValue},
		{input: false, expected: util.FalseValue},
		{input: 42, expected: "42"},
		{input: int8(42), expected: "42"},
		{input: int16(42), expected: "42"},
		{input: int32(42), expected: "42"},
		{input: int64(42), expected: "42"},
		{input: uint(42), expected: "42"},
		{input: uint8(42), expected: "42"},
		{input: uint16(42), expected: "42"},
		{input: uint32(42), expected: "42"},
		{input: uint64(42), expected: "42"},
		{input: uintptr(42), expected: "42"},
		{input: 3.1415, expected: "3.1415"},
		{input: float32(3.1415), expected: "3.1415"},
		{input: nil, expected: "<nil>"},
	}

	for _, test := range tests {
		actual := util.ValueToString(test.input)
		assert.Equal(t, test.expected, actual)
	}
}
