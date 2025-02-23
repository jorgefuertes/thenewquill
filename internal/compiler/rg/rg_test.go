package rg_test

import (
	"regexp"
	"testing"

	"thenewquill/internal/compiler/rg"

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
			rg:          rg.Blank,
			shouldMatch: true,
		},
		{
			name:        "not blank",
			text:        `not blank`,
			rg:          rg.Blank,
			shouldMatch: false,
		},
		{
			name:        "inline comment",
			text:        `foo: bar, baz // This is an inline comment`,
			rg:          rg.InlineComment,
			shouldMatch: true,
		},
		{
			name:        "inline comment 2",
			text:        `foo: bar, baz /* This is an inline comment */`,
			rg:          rg.InlineComment,
			shouldMatch: true,
		},
		{
			name:        "one line comment",
			text:        `// This is an inline comment`,
			rg:          rg.InlineComment,
			shouldMatch: true,
		},
		{
			name:        "one line comment 2",
			text:        `/* This is an inline comment */`,
			rg:          rg.InlineComment,
			shouldMatch: true,
		},
		{
			name:        "comment begin",
			text:        `/* Comment begin`,
			rg:          rg.CommentBegin,
			shouldMatch: true,
		},
		{
			name:        "comment end",
			text:        `end of the comment */`,
			rg:          rg.CommentEnd,
			shouldMatch: true,
		},
		{
			name:        "word",
			text:        `foo: bar`,
			rg:          rg.Word,
			shouldMatch: true,
		},
		{
			name:        "word with synonyms",
			text:        `foo: bar, baz, qux`,
			rg:          rg.Word,
			shouldMatch: true,
		},
		{
			name:        "include",
			text:        `INCLUDE "foo.bar"`,
			rg:          rg.Include,
			shouldMatch: true,
			matches:     []string{"foo.bar"},
		},
		{
			name:        "bad include",
			text:        `INCLUDE foo.bar`,
			rg:          rg.Include,
			shouldMatch: false,
		},
		{
			name:        "section",
			text:        `seCtion vars`,
			rg:          rg.Section,
			shouldMatch: true,
			matches:     []string{"vars"},
		},
		{
			name:        "bad section",
			text:        `bad section vars`,
			rg:          rg.Section,
			shouldMatch: false,
		},
		{
			name:        "var declaration",
			text:        `foo = "bar"`,
			rg:          rg.Var,
			shouldMatch: true,
			matches:     []string{"foo", "bar"},
		},
		{
			name:        "float",
			text:        `0.256`,
			rg:          rg.Float,
			shouldMatch: true,
		},
		{
			name:        "int",
			text:        `256`,
			rg:          rg.Int,
			shouldMatch: true,
		},
		{
			name:        "bool true",
			text:        `true`,
			rg:          rg.Bool,
			shouldMatch: true,
		},
		{
			name:        "bool false",
			text:        `false`,
			rg:          rg.Bool,
			shouldMatch: true,
		},
		{
			name:        "message",
			text:        `foo: "This is a message"`,
			rg:          rg.Msg,
			shouldMatch: true,
			matches:     []string{"foo", "This is a message"},
		},
		{
			name:        "location conns",
			text:        `exits: north loc001, south loc002, east loc003, west loc004`,
			rg:          rg.LocConns,
			shouldMatch: true,
		},
		{
			name:        "location conns2",
			text:        `exits: salir vestíbulo`,
			rg:          rg.LocConns,
			shouldMatch: true,
		},
		{
			name:        "item declaration",
			text:        `antorcha: antorcha _`,
			rg:          rg.ItemDeclaration,
			shouldMatch: true,
		},
		{
			name:        "item declaration 2",
			text:        `llave: llave dorada`,
			rg:          rg.ItemDeclaration,
			shouldMatch: true,
		},
		{
			name:        "item location",
			text:        `is in loc-004`,
			rg:          rg.ItemLocation,
			shouldMatch: true,
		},
		{
			name:        "item weight",
			text:        `has weight 10`,
			rg:          rg.ItemWeight,
			shouldMatch: true,
		},
		{
			name:        "item max weight",
			text:        `has max weight 250`,
			rg:          rg.ItemMaxWeight,
			shouldMatch: true,
		},
		{
			name:        "pluralized message",
			text:        `foo.zero: "No foos."`,
			rg:          rg.Msg,
			shouldMatch: true,
			matches:     []string{"foo.zero", "No foos."},
		},
		{
			name:        "pluralized message",
			text:        `foo.one: "One foo."`,
			rg:          rg.Msg,
			shouldMatch: true,
			matches:     []string{"foo.one", "One foo."},
		},
		{
			name:        "pluralized message",
			text:        `foo.more: "We have _ foos."`,
			rg:          rg.MsgPlural,
			shouldMatch: true,
			matches:     []string{"foo", "more", "We have _ foos."},
		},
		{
			name:        "not pluralized message",
			text:        `foo.bar: "We have _ foos."`,
			rg:          rg.MsgPlural,
			shouldMatch: false,
		},
		{
			name:        "not pluralized message",
			text:        `foo.bar: "We have _ foos."`,
			rg:          rg.Msg,
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
