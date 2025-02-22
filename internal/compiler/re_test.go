package compiler

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegexps(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		rg          *regexp.Regexp
		shouldMatch bool
		matches     []string
	}{
		{
			name:        "blank",
			text:        "    \t",
			rg:          blankRg,
			shouldMatch: true,
		},
		{
			name:        "not blank",
			text:        `not blank`,
			rg:          blankRg,
			shouldMatch: false,
		},
		{
			name:        "inline comment",
			text:        `foo: bar, baz // This is an inline comment`,
			rg:          inlineCommentRg,
			shouldMatch: true,
		},
		{
			name:        "inline comment 2",
			text:        `foo: bar, baz /* This is an inline comment */`,
			rg:          inlineCommentRg,
			shouldMatch: true,
		},
		{
			name:        "one line comment",
			text:        `// This is an inline comment`,
			rg:          inlineCommentRg,
			shouldMatch: true,
		},
		{
			name:        "one line comment 2",
			text:        `/* This is an inline comment */`,
			rg:          inlineCommentRg,
			shouldMatch: true,
		},
		{
			name:        "comment begin",
			text:        `/* Comment begin`,
			rg:          commentBeginRg,
			shouldMatch: true,
		},
		{
			name:        "comment end",
			text:        `end of the comment */`,
			rg:          commentEndRg,
			shouldMatch: true,
		},
		{
			name:        "word",
			text:        `foo: bar`,
			rg:          wordRg,
			shouldMatch: true,
		},
		{
			name:        "word with synonyms",
			text:        `foo: bar, baz, qux`,
			rg:          wordRg,
			shouldMatch: true,
		},
		{
			name:        "include",
			text:        `INCLUDE "foo.bar"`,
			rg:          includeRg,
			shouldMatch: true,
			matches:     []string{"foo.bar"},
		},
		{
			name:        "bad include",
			text:        `INCLUDE foo.bar`,
			rg:          includeRg,
			shouldMatch: false,
		},
		{
			name:        "section",
			text:        `seCtion vars`,
			rg:          sectionRg,
			shouldMatch: true,
			matches:     []string{"vars"},
		},
		{
			name:        "bad section",
			text:        `bad section vars`,
			rg:          sectionRg,
			shouldMatch: false,
		},
		{
			name:        "var declaration",
			text:        `foo = "bar"`,
			rg:          varRg,
			shouldMatch: true,
			matches:     []string{"foo", "bar"},
		},
		{
			name:        "float",
			text:        `0.256`,
			rg:          floatRg,
			shouldMatch: true,
		},
		{
			name:        "int",
			text:        `256`,
			rg:          intRg,
			shouldMatch: true,
		},
		{
			name:        "bool true",
			text:        `true`,
			rg:          boolRg,
			shouldMatch: true,
		},
		{
			name:        "bool false",
			text:        `false`,
			rg:          boolRg,
			shouldMatch: true,
		},
		{
			name:        "message",
			text:        `foo: "This is a message"`,
			rg:          msgRg,
			shouldMatch: true,
			matches:     []string{"foo", "This is a message"},
		},
		{
			name:        "location conns",
			text:        `exits: north loc001, south loc002, east loc003, west loc004`,
			rg:          locConnsRg,
			shouldMatch: true,
		},
		{
			name:        "location conns2",
			text:        `exits: salir vestíbulo`,
			rg:          locConnsRg,
			shouldMatch: true,
		},
		{
			name:        "item declaration",
			text:        `antorcha: antorcha _`,
			rg:          itemDeclarationRg,
			shouldMatch: true,
		},
		{
			name:        "item declaration 2",
			text:        `llave: llave dorada`,
			rg:          itemDeclarationRg,
			shouldMatch: true,
		},
		{
			name:        "item location",
			text:        `is in loc-004`,
			rg:          itemLocationRg,
			shouldMatch: true,
		},
		{
			name:        "item weight",
			text:        `has weight 10`,
			rg:          itemWeightRg,
			shouldMatch: true,
		},
		{
			name:        "item max weight",
			text:        `has max weight 250`,
			rg:          itemMaxWeightRg,
			shouldMatch: true,
		},
		{
			name:        "pluralized message",
			text:        `foo.zero: "No foos."`,
			rg:          msgRg,
			shouldMatch: true,
			matches:     []string{"foo.zero", "No foos."},
		},
		{
			name:        "pluralized message",
			text:        `foo.one: "One foo."`,
			rg:          msgRg,
			shouldMatch: true,
			matches:     []string{"foo.one", "One foo."},
		},
		{
			name:        "pluralized message",
			text:        `foo.more: "We have _ foos."`,
			rg:          msgPluralRg,
			shouldMatch: true,
			matches:     []string{"foo", "more", "We have _ foos."},
		},
		{
			name:        "not pluralized message",
			text:        `foo.bar: "We have _ foos."`,
			rg:          msgPluralRg,
			shouldMatch: false,
		},
		{
			name:        "not pluralized message",
			text:        `foo.bar: "We have _ foos."`,
			rg:          msgRg,
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.shouldMatch {
				require.False(t, tt.rg.MatchString(tt.text), "matches for: %v in '%s'", tt.rg, tt.text)

				return
			}

			require.True(t, tt.rg.MatchString(tt.text), "no matches for: %v in '%s'", tt.rg, tt.text)
			matches := tt.rg.FindStringSubmatch(tt.text)
			if len(tt.matches) > 0 {
				require.Equal(
					t,
					len(tt.matches),
					len(matches)-1,
					"%d expected matches, got %d: %v",
					len(tt.matches),
					len(matches)-1,
					matches,
				)

				assert.Equal(t, tt.text, matches[0])
				for i, expected := range tt.matches {
					assert.Equal(t, expected, matches[i+1])
				}
			}
		})
	}
}

func TestLabelAndTextRg(t *testing.T) {
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
			l := line{text: tt.lineText, n: 0}
			result, ok := l.labelAndTextRg(tt.label)
			require.Equal(t, tt.shouldMatch, ok)
			assert.Equal(t, tt.expected, result)
		})
	}
}
