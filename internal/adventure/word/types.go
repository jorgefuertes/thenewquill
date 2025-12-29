package word

type WordType byte

const (
	None WordType = iota
	Verb
	Noun
	Pronoun
	Adjective
	Adverb
	Preposition
	Conjunction
)

var WordTypes = []WordType{
	Verb,
	Noun,
	Pronoun,
	Adjective,
	Adverb,
	Preposition,
	Conjunction,
}

func (w WordType) String() string {
	switch w {
	case None:
		return "none"
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
		return "none"
	}
}

func WordTypeFromString(s string) WordType {
	switch s {
	case "none":
		return None
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
		return None
	}
}

func (w WordType) Byte() byte {
	return byte(w)
}

func WordTypeFromByte(b byte) WordType {
	return WordType(b)
}
