package log

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarningLevel
	ErrorLevel
	NoLevel
)

var (
	defaultLevel = WarningLevel
	output       io.Writer
)

func (l LogLevel) String() string {
	switch l {
	case InfoLevel:
		return "INFO"
	case WarningLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		return "DEBUG"
	}
}

func isTerminal() bool {
	return output != nil && output == os.Stdout
}

func (l LogLevel) Color() string {
	if !isTerminal() {
		return fmt.Sprintf("[%5s]", l.String())
	}

	var c color.Attribute
	switch l {
	case DebugLevel:
		c = color.FgHiBlue
	case InfoLevel:
		c = color.FgHiCyan
	case WarningLevel:
		c = color.FgHiYellow
	case ErrorLevel:
		c = color.FgHiRed
	default:
		c = color.FgWhite
	}

	return color.New(c).Sprintf("[%5s]", l)
}

func SetLevel(l LogLevel) {
	defaultLevel = l
}

func SetOutput(w io.Writer) {
	if w == nil {
		output = os.Stdout

		return
	}

	output = w
}

func send(level LogLevel, format string, args ...any) {
	if output == nil {
		output = os.Stdout
	}

	if defaultLevel > level {
		return
	}

	for _, line := range getLines(format, args...) {
		if level == NoLevel {
			if _, err := fmt.Fprintln(output, line); err != nil {
				fmt.Println(line)
				fmt.Printf("[ERROR] I can't write to the output: %s\n", err)
			}

			continue
		}

		if _, err := fmt.Fprintf(output, "%s %s\n", level.Color(), line); err != nil {
			fmt.Printf("[%5s] %s\n", level.String(), line)
			fmt.Printf("[ERROR] I can't write to the output: %s", err)
		}
	}
}

func getLines(format string, args ...any) []string {
	if len(args) == 0 {
		return strings.Split(format, "\n")
	}

	return strings.Split(fmt.Sprintf(format, args...), "\n")
}

func Debug(format string, args ...any) {
	send(DebugLevel, format, args...)
}

func Info(format string, args ...any) {
	send(InfoLevel, format, args...)
}

func Warning(format string, args ...any) {
	send(WarningLevel, format, args...)
}

func Error(format string, args ...any) {
	send(ErrorLevel, format, args...)
}

func Fatal(format string, args ...any) {
	send(ErrorLevel, format, args...)
	os.Exit(1)
}

func WithoutLevel(format string, args ...any) {
	send(NoLevel, format, args...)
}

func WithoutFormat(level LogLevel, line string) {
	send(level, "%s", line)
}
