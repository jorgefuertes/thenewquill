package line_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/compiler/line"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetIndent(t *testing.T) {
	m := line.NewMulti(
		line.New(`test: """`, 1),
		line.New("    This is a test", 2),
		line.New("        This is indented", 3),
		line.New(`"""`, 4),
	)

	expectedIndent := "    "
	assert.Equal(t, expectedIndent, m.GetIndent(), "Indentation should match expected value")
}

func TestStart(t *testing.T) {
	m := line.NewMulti(line.New("First line", 0))
	lineB := line.New("Another test line", 1)
	m.Append(lineB)

	assert.True(t, m.IsOn(), "Multi line should be on after starting it")
	require.Equal(t, 2, m.Len(), "There should be exactly 1 line in multiLine")
	l, ok := m.GetByIndex(1)
	require.True(t, ok)
	assert.Equal(t, lineB, l, "The initial line should match")
}

func TestJoin(t *testing.T) {
	testCases := []struct {
		name     string
		lineText []string
		expected string
	}{
		{
			name: "tabbed end",
			lineText: []string{
				`multi: """`,
				`    This is a test line \`,
				`    that continues`,
				`    	with this indented`,
				`    and ends here.`,
				`	"""`,
			},
			expected: "multi: \"This is a test line that continues\n\twith this indented\nand ends here.\"",
		},
		{
			name: "normal end",
			lineText: []string{
				`multi: """`,
				`    This is a test line \`,
				`    that continues`,
				`    	with this indented`,
				`    and ends here.`,
				`"""`,
			},
			expected: "multi: \"This is a test line that continues\n\twith this indented\nand ends here.\"",
		},
		{
			name: "spaced end",
			lineText: []string{
				`multi: """`,
				`    This is a test line \`,
				`    that continues`,
				`    	with this indented`,
				`    and ends here.`,
				`  """  	`,
			},
			expected: "multi: \"This is a test line that continues\n\twith this indented\nand ends here.\"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			m := line.NewMulti()
			for i, text := range tt.lineText {
				l := line.New(text, i)
				if l.Text() == `"""` {
					require.True(t, l.IsMultilineEnd(m.IsHeredoc()))
				}

				m.Append(l)
			}

			joinedLine := m.Join()

			assert.Equal(t, tt.expected, joinedLine.Text())
		})
	}
}

func TestJoinBackslashed(t *testing.T) {
	linesText := []string{
		`multi: "This is a test line \`,
		`	that continues \`,
		`	with no indent at all \`,
		`	and ends here."`,
	}

	m := line.NewMulti()

	for i, l := range linesText {
		m.Append(line.New(l, i))
	}

	joinedLine := m.Join()
	expectedText := "multi: \"This is a test line that continues with no indent at all and ends here.\""

	assert.True(t, m.IsOn())
	assert.Equal(t, expectedText, joinedLine.Text())

	m.Clear()
	assert.False(t, m.IsOn())
	require.Equal(t, 0, m.Len())
}
