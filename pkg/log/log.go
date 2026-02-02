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
	FatalLevel
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
	case FatalLevel:
		return "FATAL"
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
			if _, err := fmt.Fprintln(output, line); err != nil {
				fmt.Println(line)
				fmt.Printf("[ERROR] I can't write to the output: %s\n", err)
			}

			continue
		}

		if _, err := fmt.Fprintf(output, "%s %s\n", level.color(), line); err != nil {
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
