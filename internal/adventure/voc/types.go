package voc

import "strings"

type WordType int

const (
	Verb WordType = iota
	Noun
	Pronoun
	Adjective
	Adverb
	Preposition
	Conjunction
	Unknown
)

func wordTypes() []WordType {
	return []WordType{Verb, Noun, Pronoun, Adjective, Adverb, Preposition, Conjunction}
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
