package util

import (
	"fmt"
	"strconv"
	"strings"
)

func LimitStr(s string, max int) string {
	runes := []rune(s)

	if len(runes) >= max-3 && len(runes) > 15 {
		return string(runes[:max-3]) + "..."
	}

	if len(runes) > max {
		return string(runes[:max])
	}

	return string(runes)
}

func RemoveAccents(s string) string {
	replacements := map[string]string{
		`á`: `a`, `é`: `e`, `í`: `i`, `ó`: `o`, `ú`: `u`,
		`ñ`: `n`, `ü`: `u`, `ç`: `c`,
	}

	for old, new := range replacements {
		s = strings.ReplaceAll(s, old, new)
	}

	return s
}

func RemoveSymbols(s string) string {
	var output string

	for i := 0; i < len(s); i++ {
		if !strings.Contains("!@#$%^&*()-+={}[]:;\"'<>,.?/|\\", string(s[i])) {
			output += string(s[i])
		}
	}

	return output
}

func ValueToString(v any) string {
	switch v.(type) {
	case string:
		return fmt.Sprintf("s:%s", v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("i:%d", v)
	case float32, float64:
		return fmt.Sprintf("f:%.4f", v)
	case bool:
		if v == true {
			return "b:T"
		} else {
			return "b:F"
		}
	}

	return fmt.Sprintf("s:%v", v)
}

func StringToValue(s string) any {
	switch s[0] {
	case 's':
		return s[2:]
	case 'i':
		i, _ := strconv.Atoi(s[2:])
		return i
	case 'f':
		f, _ := strconv.ParseFloat(s[2:], 64)
		return f
	case 'b':
		if s[2:] == "T" {
			return true
		} else {
			return false
		}
	}

	panic("invalid string")
}

func SplitString(text string, chunkSize int) []string {
	var chunks []string
	for i := 0; i < len(text); i += chunkSize {
		end := i + chunkSize
		if end > len(text) {
			end = len(text)
		}
		chunks = append(chunks, text[i:end])
	}

	return chunks
}

func EscapeExportString(s string) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	s = strings.ReplaceAll(s, "|", "\\|")

	return s
}

func UnescapeExportString(s string) string {
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.ReplaceAll(s, "\\r", "\r")
	s = strings.ReplaceAll(s, "\\t", "\t")
	s = strings.ReplaceAll(s, "\\|", "|")

	return s
}
