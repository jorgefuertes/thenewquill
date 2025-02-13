package voc

import "slices"

type Vocabulary []*Word

func New() Vocabulary {
	v := Vocabulary{}

	for _, t := range wordTypes() {
		v = append(v, &Word{Label: "_", Type: t, Synonyms: make([]string, 0)})
	}

	return v
}

func (v *Vocabulary) Add(label string, t WordType, synonyms ...string) *Word {
	if synonyms == nil {
		synonyms = make([]string, 0)
	}

	if v.Exists(t, label) {
		w := v.Get(t, label)
		w.Synonyms = synonyms

		if w.Type == Unknown {
			w.Type = t
		}

		return w
	}

	w := &Word{
		Label:    label,
		Type:     t,
		Synonyms: synonyms,
	}

	*v = append(*v, w)

	return w
}

func (v Vocabulary) Get(t WordType, labelOrSynonym string) *Word {
	for _, w := range v {
		if w.Is(labelOrSynonym) && w.Type == t {
			return w
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

// First returns the first word of the vocabulary that matches the given label or synonym.
func (v Vocabulary) First(labelOrSynonym string) *Word {
	for _, w := range v {
		if w.Is(labelOrSynonym) {
			return w
		}
	}

	return nil
}

// FirstWithTypes returns the first word found by given types preferences
func (v Vocabulary) FirstWithTypes(labelOrSynonym string, types ...WordType) *Word {
	for _, t := range types {
		w := v.Get(t, labelOrSynonym)
		if w != nil {
			return w
		}
	}

	return nil
}
