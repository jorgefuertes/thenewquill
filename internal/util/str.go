package util

import (
	"encoding/base64"
	"regexp"
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

func SplitIntoLines(text string, maxLen int) []string {
	lines := make([]string, 0)
	words := regexp.MustCompile(`\s+`).Split(text, -1)

	var line string
	for _, w := range words {
		line = strings.TrimSpace(line)

		if len(w) > maxLen {
			if line != "" {
				lines = append(lines, line)
				line = ""
			}

			lines = append(lines, SplitWithDashes(w, maxLen)...)

			continue
		}

		if len(line)+len(w)+1 > maxLen {
			if line != "" {
				lines = append(lines, line)
				line = ""
			}
		}

		if line != "" {
			line += " "
		}

		line += w
	}

	if line != "" {
		lines = append(lines, line)
	}

	return lines
}

func SplitWithDashes(text string, maxLen int) []string {
	lines := make([]string, 0)

	var line string
	for _, c := range text {
		if len(line)+1 == maxLen {
			lines = append(lines, strings.TrimSpace(line)+"-")
			line = ""

			continue
		}

		line += string(c)
	}

	line = strings.TrimSpace(line)
	if line != "" {
		lines = append(lines, line)
	}

	return lines
}

func SplitIntoFields(s string) []string {
	fields := strings.Split(s, "|")

	for i, field := range fields {
		if strings.HasPrefix(field, "@B64:") && len(field) > 5 {
			b, err := base64.StdEncoding.DecodeString(field[4:])
			if err == nil {
				fields[i] = string(b)
			}
		}
	}

	return fields
}

func EscapeField(s string) string {
	return "@B64:" + base64.StdEncoding.EncodeToString([]byte(s))
}

func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)

	return i
}
