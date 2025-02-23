package line_test

import (
	"testing"

	"thenewquill/internal/compiler/line"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTextForLabel(t *testing.T) {
	tests := []struct {
		name        string
		lineText    string
		label       string
		expected    string
		shouldMatch bool
	}{
		{
			name:        "title",
			lineText:    `title: "Catacombs"`,
			label:       "title",
			expected:    `Catacombs`,
			shouldMatch: true,
		},
		{
			name:        "title with weird spacing",
			lineText:    `	 title: 	"Catacombs" `,
			label:       "title",
			expected:    `Catacombs`,
			shouldMatch: true,
		},
		{
			name:        "title with comment",
			lineText:    `title: "Catacombs" // comment`,
			label:       "title",
			expected:    `Catacombs`,
			shouldMatch: true,
		},
		{
			name:        "desc",
			lineText:    `desc: "In a dark cave, you see several niches and a large chamber."`,
			label:       "desc",
			expected:    `In a dark cave, you see several niches and a large chamber.`,
			shouldMatch: true,
		},
		{
			name:        "desc with weird spacing",
			lineText:    `desc:   "In a dark cave, you see several niches and a large chamber."	    `,
			label:       "desc",
			expected:    `In a dark cave, you see several niches and a large chamber.`,
			shouldMatch: true,
		},
		{
			name:        "desc with colons",
			lineText:    `desc: "In a \"dark cave\", you see several niches and a large 'chamber'."`,
			label:       "desc",
			expected:    `In a "dark cave", you see several niches and a large 'chamber'.`,
			shouldMatch: true,
		},
		{
			name:        "no match",
			lineText:    `foo: "No Match"`,
			label:       "bar",
			expected:    "",
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := line.New(tt.lineText, 0)
			result, ok := l.GetTextForLabel(tt.label)
			require.Equal(t, tt.shouldMatch, ok)
			assert.Equal(t, tt.expected, result)
		})
	}
}
