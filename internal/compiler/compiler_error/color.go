package compiler_error

import (
	"regexp"

	"github.com/charmbracelet/lipgloss"
)

var (
	sectionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ddff")).Bold(true)
	pathStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#eeaa00"))
	numberStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#666")).Background(lipgloss.Color("#000099"))
	redStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#dd0000"))
	borderStyle  = lipgloss.NewStyle().
			Border(lipgloss.BlockBorder()).
			BorderForeground(lipgloss.Color("#CC3322")).
			Padding(1, 1, 0, 1)
	errorTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#DD2222")).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)
	directiveStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#11ccee"))
	labelStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#880088"))
	reservedWordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#2277cc"))
	labelNumberStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#adad22"))
)

const (
	upArrowHead   = "\U000002C4"
	downArrowHead = "\U000002C5"
)

type decoration struct {
	rg     string
	style  lipgloss.Style
	render string
}

var decorators = []decoration{
	{`\b(\d{5})\b`, numberStyle, "$1"},                                      // line numbers
	{`(\#\d+)`, labelNumberStyle, "$1"},                                     // label numbers
	{`(\x{2C4}|\x{2C5})`, redStyle.Bold(true), "$1"},                        // arrows
	{`\b([\p{L}\.\/]+\.\p{L}{3,6})\b`, pathStyle, "$1"},                     // files
	{`SECTION ([A-Z\s]{3,20}){1,3}`, sectionStyle.Bold(true), "SECTION $1"}, // section
	{`\b(FILE|INCLUDE)\b`, directiveStyle, "$1"},                            // compiler directives
	{`(?i)\b(id|label)\b`, reservedWordStyle, "$1"},                         // reserved words
	{`\b(ERROR|FATAL|PANIC)\b`, redStyle.Bold(true), "$1"},                  // critical error
	{`\b([A-Za-z0-9_-]{3,25})\:(\s+)`, labelStyle, `$1:$2`},                 // label
}

func decorateLine(line string) string {
	for _, dec := range decorators {
		rg := regexp.MustCompile(dec.rg)
		if rg.MatchString(line) {
			line = rg.ReplaceAllString(line, dec.style.Render(dec.render))
		}
	}

	return line
}
