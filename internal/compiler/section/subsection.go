package section

import (
	"strings"
)

type SubSection byte

const (
	NoSubSection SubSection = 0
	Verb         SubSection = 1
	Noun         SubSection = 2
	Pronoun      SubSection = 3
	Adjective    SubSection = 4
	Adverb       SubSection = 5
	Preposition  SubSection = 6
	Conjunction  SubSection = 7
)

func SubSections() []SubSection {
	return []SubSection{NoSubSection, Verb, Noun, Pronoun, Adjective, Adverb, Preposition, Conjunction}
}

func SubSectionNames() map[SubSection]string {
	return map[SubSection]string{
		NoSubSection: "none",
		Verb:         "verb",
		Noun:         "noun",
		Pronoun:      "pronoun",
		Adjective:    "adjective",
		Adverb:       "adverb",
		Preposition:  "preposition",
		Conjunction:  "conjunction",
	}
}

func (s SubSection) Byte() byte {
	return byte(s)
}

func SubSectionFromByte(b byte) SubSection {
	if int(b) < 0 || int(b) >= len(SubSections()) {
		return NoSubSection
	}

	return SubSection(b)
}

func SubSectionFromString(s string) SubSection {
	s = strings.ToLower(s)

	for sub, name := range SubSectionNames() {
		if name == s {
			return sub
		}
	}

	return NoSubSection
}
