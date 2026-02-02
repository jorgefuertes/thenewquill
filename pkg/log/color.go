package log

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func isTerminal() bool {
	return output != nil && output == os.Stdout
}

func (level LogLevel) color() string {
	if !isTerminal() || termenv.ColorProfile() == termenv.Ascii {
		return fmt.Sprintf("[%s]", level.String())
	}

	darkBgColors := map[LogLevel]lipgloss.Style{
		DebugLevel:   lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")),
		InfoLevel:    lipgloss.NewStyle().Foreground(lipgloss.Color("#00ffff")),
		WarningLevel: lipgloss.NewStyle().Foreground(lipgloss.Color("#ffff00")),
		ErrorLevel:   lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6666")),
		FatalLevel:   lipgloss.NewStyle().Foreground(lipgloss.Color("#ff4444")).Bold(true),
		NoLevel:      lipgloss.NewStyle().Foreground(lipgloss.Color("#cccccc")),
	}

	lightBgColors := map[LogLevel]lipgloss.Style{
		DebugLevel:   lipgloss.NewStyle().Foreground(lipgloss.Color("#006600")),
		InfoLevel:    lipgloss.NewStyle().Foreground(lipgloss.Color("#008888")),
		WarningLevel: lipgloss.NewStyle().Foreground(lipgloss.Color("#888800")),
		ErrorLevel:   lipgloss.NewStyle().Foreground(lipgloss.Color("#cc4444")),
		FatalLevel:   lipgloss.NewStyle().Foreground(lipgloss.Color("#aa0000")).Bold(true),
		NoLevel:      lipgloss.NewStyle().Foreground(lipgloss.Color("#333333")),
	}

	if termenv.HasDarkBackground() {
		if style, ok := darkBgColors[level]; ok {
			return fmt.Sprintf("%s", style.Render(level.String()))
		}

		return fmt.Sprintf("%s", darkBgColors[NoLevel].Render(NoLevel.String()))
	}

	if style, ok := lightBgColors[level]; ok {
		return fmt.Sprintf("%s", style.Render(level.String()))
	}

	return fmt.Sprintf("%s", lightBgColors[NoLevel].Render(NoLevel.String()))
}
