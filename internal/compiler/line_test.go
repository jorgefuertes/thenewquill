package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLineToVar(t *testing.T) {
	type result struct {
		key   string
		value any
		ok    bool
	}

	testCases := []struct {
		name     string
		lineText string
		expected result
	}{
		{
			name:     "string declaration",
			lineText: `test = "This is a test string var"`,
			expected: result{
				key:   "test",
				value: "This is a test string var",
				ok:    true,
			},
		},
		{
			name:     "with indentation",
			lineText: `	test = "This is a test string var"`,
			expected: result{
				key:   "test",
				value: "This is a test string var",
				ok:    true,
			},
		},
		{
			name:     "weird spacing",
			lineText: `	test  = 	"This is a test string var"  `,
			expected: result{
				key:   "test",
				value: "This is a test string var",
				ok:    true,
			},
		},
		{
			name:     "line comment",
			lineText: `	test  = 	"This is a test string var" // comment`,
			expected: result{
				key:   "test",
				value: "This is a test string var",
				ok:    true,
			},
		},
		{
			name:     "unclosed declaration",
			lineText: `	test  = 	"This is a test string var`,
			expected: result{
				key:   "test",
				value: "This is a test string var",
				ok:    true,
			},
		},
		{
			name:     "int declaration",
			lineText: `test = 1`,
			expected: result{
				key:   "test",
				value: 1,
				ok:    true,
			},
		},
		{
			name:     "float declaration",
			lineText: `test = 0.536`,
			expected: result{
				key:   "test",
				value: 0.536,
				ok:    true,
			},
		},
		{
			name:     "bool true declaration",
			lineText: `test = true`,
			expected: result{
				key:   "test",
				value: true,
				ok:    true,
			},
		},
		{
			name:     "bool false declaration",
			lineText: `test = 	false // this is false`,
			expected: result{
				key:   "test",
				value: false,
				ok:    true,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			l := line{text: tt.lineText, n: 1}
			key, res, ok := l.toVar()
			assert.Equal(t, tt.expected.ok, ok)
			assert.Equal(t, tt.expected.key, key)
			require.EqualValues(t, tt.expected.value, res)
		})
	}
}
