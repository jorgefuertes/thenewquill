package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetIndent(t *testing.T) {
	s := newStatus()

	lines := []line{
		{text: `test: """`},
		{text: "    This is a test"},
		{text: "        This is indented"},
		{text: `"""`},
	}
	s.multiLine = multi{lines: lines}

	expectedIndent := "    "
	assert.Equal(t, expectedIndent, s.multiLine.getIndent(), "Indentation should match expected value")
}

func TestStartMultiLine(t *testing.T) {
	s := newStatus()

	lineB := line{text: "Another test line"}

	s.appendMultiLine(lineB)

	assert.True(t, s.multiLine.isOn(), "Multi line should be on after starting it")
	require.Equal(t, 1, len(s.multiLine.lines), "There should be exactly 1 line in multiLine")
	assert.Equal(t, lineB, s.multiLine.lines[0], "The initial line should match expected")
}

func TestJoinAndClearMultiLine(t *testing.T) {
	s := newStatus()

	lines := []line{
		{text: `multi: """`},
		{text: `    This is a test line \`},
		{text: `    that continues`},
		{text: `    	with this indented`},
		{text: `    and ends here.`},
		{text: `"""`},
	}

	for _, l := range lines {
		if l.text == `"""` {
			require.True(t, l.isMultilineEnd(s.multiLine.isHeredoc()))
		}

		s.appendMultiLine(l)
	}

	joinedLine := s.joinAnClearMultiLine()
	expectedText := "multi: \"This is a test line that continues\n\twith this indented\nand ends here.\""

	assert.Equal(t, expectedText, joinedLine.text)
	assert.False(t, s.multiLine.isOn())
	require.Equal(t, 0, len(s.multiLine.lines))
}

func TestJoinAndClearMultiLineBackslashed(t *testing.T) {
	s := newStatus()

	lines := []line{
		{text: `multi: "This is a test line \`},
		{text: `	that continues \`},
		{text: `	with no indent at all \`},
		{text: `	and ends here."`},
	}

	for _, l := range lines {
		s.appendMultiLine(l)
	}

	joinedLine := s.joinAnClearMultiLine()
	expectedText := "multi: \"This is a test line that continues with no indent at all and ends here.\""

	assert.Equal(t, expectedText, joinedLine.text)
	assert.False(t, s.multiLine.isOn())
	require.Equal(t, 0, len(s.multiLine.lines))
}
