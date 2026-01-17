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

func GetDefaultSynonymForAction(lang Lang, action Action) []string {
	switch action {
	case Go:
		if lang == ES {
			return []string{"ir", "moverse", "caminar", "ve", "muévete", "camina"}
		}

		return []string{"go", "move", "walk"}
	case Examine:
		if lang == ES {
			return []string{"ex", "examinar", "mirar", "ver", "examina", "mira"}
		}

		return []string{"ex", "examine", "look", "see"}
	case Talk:
		if lang == ES {
			return []string{"hablar", "decir", "charlar", "di", "habla", "charla"}
		}

		return []string{"talk", "say", "chat"}
	default:
		return []string{}
	}
}
