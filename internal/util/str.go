package util

import "strings"

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
