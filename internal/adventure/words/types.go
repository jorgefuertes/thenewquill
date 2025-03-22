package words

import "strings"

type WordType byte

const (
	Unknown     WordType = 0
	Verb        WordType = 1
	Noun        WordType = 2
	Pronoun     WordType = 3
	Adjective   WordType = 4
	Adverb      WordType = 5
	Preposition WordType = 6
	Conjunction WordType = 7
)

func wordTypes() []WordType {
	return []WordType{Unknown, Verb, Noun, Pronoun, Adjective, Adverb, Preposition, Conjunction}
}

func (t WordType) String() string {
	switch t {
	case Verb:
		return "verb"
	case Noun:
		return "noun"
	case Pronoun:
		return "pronoun"
	case Adjective:
		return "adjective"
	case Adverb:
		return "adverb"
	case Preposition:
		return "preposition"
	case Conjunction:
		return "conjunction"
	default:
		return "unknown"
	}
}

func (t WordType) Int() int {
	return int(t)
}

func WordTypeFromString(s string) WordType {
	s = strings.ToLower(s)

	switch s {
	case "verb":
		return Verb
	case "noun":
		return Noun
	case "pronoun":
		return Pronoun
	case "adjective":
		return Adjective
	case "adverb":
		return Adverb
	case "preposition":
		return Preposition
	case "conjunction":
		return Conjunction
	default:
		return Unknown
	}
}
