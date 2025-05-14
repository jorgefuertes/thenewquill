package db

import (
	"strings"
)

type SubKind byte

const (
	AnySubKind         SubKind = 255
	NoSubKind          SubKind = 0
	VerbSubKind        SubKind = 1
	NounSubKind        SubKind = 2
	PronounSubKind     SubKind = 3
	AdjectiveSubKind   SubKind = 4
	AdverbSubKind      SubKind = 5
	PrepositionSubKind SubKind = 6
	ConjunctionSubKind SubKind = 7
)

func (k Kind) SubKinds() []SubKind {
	switch k {
	case Words:
		return []SubKind{
			VerbSubKind,
			NounSubKind,
			PronounSubKind,
			AdjectiveSubKind,
			AdverbSubKind,
			PrepositionSubKind,
			ConjunctionSubKind,
		}
	default:
		return []SubKind{NoSubKind}
	}
}

func SubKinds() []SubKind {
	return []SubKind{
		NoSubKind,
		VerbSubKind,
		NounSubKind,
		PronounSubKind,
		AdjectiveSubKind,
		AdverbSubKind,
		PrepositionSubKind,
		ConjunctionSubKind,
	}
}

func SubKindNames() map[SubKind]string {
	return map[SubKind]string{
		NoSubKind:          "none",
		VerbSubKind:        "verb",
		NounSubKind:        "noun",
		PronounSubKind:     "pronoun",
		AdjectiveSubKind:   "adjective",
		AdverbSubKind:      "adverb",
		PrepositionSubKind: "preposition",
		ConjunctionSubKind: "conjunction",
	}
}

func (s SubKind) String() string {
	name, ok := SubKindNames()[s]
	if !ok {
		return SubKindNames()[NoSubKind]
	}

	return name
}

func (s SubKind) Byte() byte {
	return byte(s)
}

func SubKindFromByte(b byte) SubKind {
	if int(b) < 0 || int(b) >= len(SubKinds()) {
		return NoSubKind
	}

	return SubKind(b)
}

func SubKindFromString(s string) SubKind {
	s = strings.ToLower(s)

	for sub, name := range SubKindNames() {
		if name == s {
			return sub
		}
	}

	return NoSubKind
}
