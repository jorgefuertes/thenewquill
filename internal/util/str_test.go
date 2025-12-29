package util_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestStrHelpers(t *testing.T) {
	const maxLen = 40

	text := []string{
		`123456789|123456789|123456789|123456789|`,
		`Lorem ipsum dolor sit amet, consectetur `,
		`adipiscing elit, sed do eiusmod tempor  `,
		`incididunt ut labore et dolore magna    `,
		`aliqua. Ut enim ad minim veniam, quis   `,
		`nostrud exercitation ullamco laboris    `,
		`nisi ut aliquip ex ea commodo consequat.`,
		`Duis aute irure dolor in reprehenderit  `,
		`in voluptate velit esse cillum dolore.  `,
		`eu fugiat nulla pariatur. Excepteur sint`,
		`occaecat cupidatat non proident, sunt in`,
		`in culpa qui officia deserunt mollit    `,
		`anim id est laborum.         		     `,
	}

	oneLineText := strings.Join(text, " ")
	oneLineText = regexp.MustCompile(`\s+`).ReplaceAllString(oneLineText, " ")
	oneLineText = strings.TrimSpace(oneLineText)

	t.Run("SplitIntoLines", func(t *testing.T) {
		lines := util.SplitIntoLines(oneLineText, maxLen)

		for i, line := range lines {
			t.Logf("%04d [%-40s] (%02d)", i, line, len(line))
			assert.NotEmpty(t, line, "line %d is empty", i)
			assert.LessOrEqual(t, len(line), maxLen,
				"line should be less or equal than %d chars and I have %d", maxLen, len(line))
			assert.Equal(t, strings.TrimSpace(text[i]), line)
		}

		assert.Len(t, lines, len(text), "should have %d lines and I have %d", len(text), len(lines))
	})
}
