package log

import (
	"fmt"
	"io"
	"os"
	"strings"
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
	defaultLevel = DebugLevel
	output       io.Writer
)

func (l LogLevel) String() string {
	switch l {
	case InfoLevel:
		return "INFO"
	case WarningLevel:
		return "WARN"
	case ErrorLevel:
		return "ERR"
	default:
		return "DEBUG"
	}
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
			fmt.Fprintln(output, line)

			continue
		}

		fmt.Fprintf(output, "[%5s] %s\n", level.String(), line)
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
