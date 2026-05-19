package line_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/compiler/line"

	"github.com/stretchr/testify/assert"
)

func TestIsBlank(t *testing.T) {
	testCases := []struct {
		text string
		want bool
	}{
		{"", true},
		{"   ", true},
		{"\t\t", true},
		{" \t \t ", true},
		{"a", false},
		{"  a  ", false},
		{"// just a comment", false},
	}

	for _, tc := range testCases {
		t.Run(tc.text, func(t *testing.T) {
			assert.Equal(t, tc.want, line.New(tc.text, 0).IsBlank())
		})
	}
}

func TestIsOneLineComment(t *testing.T) {
	testCases := []struct {
		text string
		want bool
	}{
		{"// a comment", true},
		{"  // indented comment", true},
		{"/* block on one line */", true},
		{"  /* indented block */  ", true},
		{"/* not closed", false},
		{"foo // trailing comment", false},
		{"plain text", false},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.text, func(t *testing.T) {
			assert.Equal(t, tc.want, line.New(tc.text, 0).IsOneLineComment())
		})
	}
}

func TestIsCommentBegin(t *testing.T) {
	testCases := []struct {
		text string
		want bool
	}{
		{"/* opens here", true},
		{"   /* indented", true},
		{"foo /* mid-line", false},
		{"// line comment", false},
		{"plain", false},
	}

	for _, tc := range testCases {
		t.Run(tc.text, func(t *testing.T) {
			assert.Equal(t, tc.want, line.New(tc.text, 0).IsCommentBegin())
		})
	}
}

func TestIsCommentEnd(t *testing.T) {
	testCases := []struct {
		text string
		want bool
	}{
		{"closes here */", true},
		{"closes here */   ", true},
		{"foo */ bar", false},
		{"no end at all", false},
	}

	for _, tc := range testCases {
		t.Run(tc.text, func(t *testing.T) {
			assert.Equal(t, tc.want, line.New(tc.text, 0).IsCommentEnd())
		})
	}
}

func TestIsMultilineBegin(t *testing.T) {
	testCases := []struct {
		name string
		text string
		want bool
	}{
		{"heredoc opener", `desc: """`, true},
		{"heredoc opener with trailing spaces", `desc: """  `, true},
		{"backslash continuation", `desc: "long stuff \`, true},
		{"plain line is not a multiline begin", `desc: "short"`, false},
		{"empty", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, line.New(tc.text, 0).IsMultilineBegin())
		})
	}
}

func TestIsMultilineEnd(t *testing.T) {
	t.Run("heredoc mode", func(t *testing.T) {
		// In heredoc mode, end is `"""` on its own line.
		assert.True(t, line.New(`"""`, 0).IsMultilineEnd(true))
		assert.True(t, line.New(`   """  `, 0).IsMultilineEnd(true))
		assert.False(t, line.New(`foo`, 0).IsMultilineEnd(true))
	})

	t.Run("backslash mode", func(t *testing.T) {
		// In backslash mode, end is any line that does NOT trail with `\`.
		assert.True(t, line.New(`final piece"`, 0).IsMultilineEnd(false))
		assert.False(t, line.New(`continues here \`, 0).IsMultilineEnd(false))
	})
}
