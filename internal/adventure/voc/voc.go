package voc

import "slices"

type Vocabulary []Word

func New() Vocabulary {
	v := Vocabulary{}

	for _, t := range wordTypes() {
		v = append(v, Word{Label: "_", Type: t})
	}

	return v
}

func (v *Vocabulary) Add(label string, t WordType, synonyms ...string) error {
	if v.Exists(t, label) {
		return ErrWordAlreadyExists
	}

	if synonyms == nil {
		synonyms = []string{}
	}

	w := Word{
		Label:    label,
		Type:     t,
		Synonyms: synonyms,
	}

	*v = append(*v, w)

	return nil
}

func (v Vocabulary) Get(t WordType, labelOrSynonym string) *Word {
	for _, w := range v {
		if w.Is(labelOrSynonym) && w.Type == t {
			return &w
		}
	}

	return nil
}

func (v Vocabulary) Exists(t WordType, labelOrSynonym string) bool {
	for _, w := range v {
		if w.Type != t {
			continue
		}

		if w.Label == labelOrSynonym || slices.Contains(w.Synonyms, labelOrSynonym) {
			return true
		}
	}

	return false
}
