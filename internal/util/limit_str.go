package util

import "unicode/utf8"

func LimitStr(s string, max int) string {
	if max >= utf8.RuneCountInString(s) {
		return s
	}

	var output string
	for i, r := range s {
		if i == max {
			break
		}

		output += string(r)
	}

	return output
}
