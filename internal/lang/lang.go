package lang

import "strings"

type (
	Lang   string
	Action uint8
)

const (
	EN Lang   = "en"
	ES Lang   = "es"
	Go Action = iota
	Examine
	Talk
)

func (l Lang) String() string {
	return string(l)
}

var allowedLanguages = []Lang{EN, ES}

func IsAllowedLanguage(l string) bool {
	l = strings.ToLower(l)

	for _, lang := range allowedLanguages {
		if lang.String() == l {
			return true
		}
	}

	return false
}

func AllowedLanguages() []Lang {
	return allowedLanguages
}
