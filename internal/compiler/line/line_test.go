package line_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/compiler/line"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAndNumber(t *testing.T) {
	l := line.New("hello", 7)
	assert.Equal(t, "hello", l.Text)
	assert.Equal(t, 7, l.Number())
}

func TestAdd(t *testing.T) {
	l := line.New("hello", 1)
	l.Add(" world")
	assert.Equal(t, "hello world", l.Text)
	assert.Equal(t, 1, l.Number(), "Number should not change on Add")

	// Appending to an empty line.
	empty := line.New("", 0)
	empty.Add("first")
	assert.Equal(t, "first", empty.Text)
}

func TestOptimizedText(t *testing.T) {
	testCases := []struct {
		name string
		text string
		want string
	}{
		{"trims surrounding whitespace", "   foo   ", "foo"},
		{"strips trailing line comment", `foo: "bar" // a note`, `foo: "bar"`},
		{"strips trailing block comment", `foo: "bar" /* a note */`, `foo: "bar"`},
		{"no comment is preserved", `foo: "bar"`, `foo: "bar"`},
		{"only whitespace becomes empty", "  \t  ", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, line.New(tc.text, 0).OptimizedText())
		})
	}
}

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
			l := line.New(tt.lineText, 1)
			key, res, ok := l.AsVar()
			assert.Equal(t, tt.expected.ok, ok)
			assert.Equal(t, tt.expected.key, key)
			require.EqualValues(t, tt.expected.value, res)
		})
	}
}
