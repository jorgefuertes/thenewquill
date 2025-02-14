package util

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
