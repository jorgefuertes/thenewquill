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
	if v.Exists(label) {
		return ErrWordAlreadyExists
	}

	w := Word{
		Label:    label,
		Type:     t,
		Synonyms: synonyms,
	}

	*v = append(*v, w)

	return nil
}

func (v Vocabulary) Get(labelOrSynonym string) *Word {
	for _, w := range v {
		if w.Is(labelOrSynonym) {
			return &w
		}
	}
	return nil
}

func (v Vocabulary) GetByType(t WordType, labelOrSynonym string) *Word {
	for _, w := range v {
		if w.Is(labelOrSynonym) && w.Type == t {
			return &w
		}
	}

	return nil
}

func (v Vocabulary) Exists(labelOrSynonym string) bool {
	for _, w := range v {
		if w.Label == labelOrSynonym || slices.Contains(w.Synonyms, labelOrSynonym) {
			return true
		}
	}

	return false
}
