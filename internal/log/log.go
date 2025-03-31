package log

import (
	"fmt"
	"os"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarningLevel
	ErrorLevel
)

var DefaultLevel = DebugLevel

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
	DefaultLevel = l
}

func Send(level LogLevel, format string, args ...any) {
	if DefaultLevel > level {
		return
	}

	fmt.Printf("[%5s] %s\n", level.String(), fmt.Sprintf(format, args...))
}

func Debug(format string, args ...any) {
	Send(DebugLevel, format, args...)
}

func Info(format string, args ...any) {
	Send(InfoLevel, format, args...)
}

func Warning(format string, args ...any) {
	Send(WarningLevel, format, args...)
}

func Error(err error) {
	Send(ErrorLevel, "%s", err.Error())
}

func Fatal(err error) {
	Error(err)
	os.Exit(1)
}
