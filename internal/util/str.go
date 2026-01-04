package util

import (
	"encoding/base64"
	"regexp"
	"strconv"
	"strings"
)

func LimitStr(s string, max int) string {
	runes := []rune(s)

	if len(runes) >= max-1 && len(runes) > 15 {
		return string(runes[:max-1]) + "…"
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

func SplitIntoLines(text string, maxLen int) []string {
	if maxLen <= 0 {
		return []string{}
	}

	reader := strings.NewReader(text)
	lines := make([]string, 0)
	line := ""

	b := make([]byte, 1)
	for {
		_, err := reader.Read(b)
		if err != nil {
			break
		}

		if string(b) == "\n" {
			lines = append(lines, line)
			line = ""

			continue
		}

		if len(line) == 0 && string(b) == " " {
			continue
		}

		if regexp.MustCompile(`[\.,;]+$`).MatchString(line) && !regexp.MustCompile(`\s+`).Match(b) {
			line += " "
			reader.UnreadByte()

			continue
		}

		line += string(b)

		if len(line) == maxLen {
			_, err := reader.Read(b)
			if err != nil {
				break
			}

			if string(b) == " " {
				lines = append(lines, line)
				line = ""

				continue
			}

			reader.UnreadByte()

			for !regexp.MustCompile(`[^\p{L}\p{N}]{1}$`).MatchString(line) {
				line = line[:len(line)-1]
				reader.UnreadByte()

				if len(line) == 0 {
					b := make([]byte, maxLen-1)

					_, err := reader.Read(b)
					if err != nil {
						panic(err)
					}

					break
				}
			}

			lines = append(lines, line)
			line = ""
		}
	}

	if len(line) > 0 {
		lines = append(lines, line)
	}

	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}

	return lines
}
