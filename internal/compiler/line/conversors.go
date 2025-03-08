package line

import (
	"strconv"
	"strings"

	"thenewquill/internal/adventure/config"
	"thenewquill/internal/adventure/msg"
	"thenewquill/internal/adventure/words"
	"thenewquill/internal/compiler/rg"
	"thenewquill/internal/compiler/section"
)

// AsInclude returns the include filename string and true if it was found
func (l Line) AsInclude() (string, bool) {
	if !rg.Include.MatchString(l.OptimizedText()) {
		return "", false
	}

	return rg.Include.FindStringSubmatch(l.OptimizedText())[1], true
}

// AsSection returns the section and true if it was found
func (l Line) AsSection() (section.Section, bool) {
	if !rg.Section.MatchString(l.OptimizedText()) {
		return section.None, false
	}

	return section.FromString(rg.Section.FindStringSubmatch(l.OptimizedText())[1]), true
}

// AsVar returns the variable name and value and true if it was found
func (l Line) AsVar() (string, any, bool) {
	o := l.OptimizedText()

	if !rg.Var.MatchString(o) {
		return "", "", false
	}

	name := rg.Var.FindStringSubmatch(o)[1]
	valueStr := rg.Var.FindStringSubmatch(o)[2]

	if rg.Float.MatchString(valueStr) {
		value, _ := strconv.ParseFloat(valueStr, 64)

		return name, value, true
	}

	if rg.Int.MatchString(valueStr) {
		value, _ := strconv.ParseInt(valueStr, 10, 64)

		return name, value, true
	}

	if rg.Bool.MatchString(valueStr) {
		value, _ := strconv.ParseBool(valueStr)

		return name, value, true
	}

	return name, valueStr, true
}

// AsWord returns the word and true if it was found
func (l Line) AsWord() (words.Word, bool) {
	o := l.OptimizedText()

	if !rg.Word.MatchString(o) {
		return words.Word{}, false
	}

	w := words.Word{}

	parts := strings.Split(o, ":")
	if len(parts) != 2 {
		return words.Word{}, false
	}

	w.Type = words.WordTypeFromString(parts[0])

	words := strings.Split(parts[1], ",")
	for i, word := range words {
		if i == 0 {
			w.Label = strings.TrimSpace(word)

			continue
		}

		w.Synonyms = append(w.Synonyms, strings.TrimSpace(word))
	}

	return w, true
}

func (l Line) AsMsg(t msg.MsgType) (*msg.Msg, bool) {
	if !rg.Msg.MatchString(l.text) {
		return nil, false
	}

	parts := rg.Msg.FindStringSubmatch(l.text)
	if len(parts) != 3 {
		return nil, false
	}

	return msg.New(t, parts[1], parts[2]), true
}

// AsLocationLabel returns the location label and true if it was found
func (l Line) AsLocationLabel() (string, bool) {
	if !rg.LocLabel.MatchString(l.OptimizedText()) {
		return "", false
	}

	return rg.LocLabel.FindStringSubmatch(l.OptimizedText())[1], true
}

// AsLocationDescription returns the location description and true if it was found
func (l Line) AsLocationDescription() (string, bool) {
	return l.GetTextForLabel("desc")
}

// AsLocationTitle returns the location title and true if it was found
func (l Line) AsLocationTitle() (string, bool) {
	return l.GetTextForLabel("title")
}

// AsLocationConns returns the location connections and true if it was found
func (l Line) AsLocationConns() (map[string]string, bool) {
	exits := make(map[string]string, 0)

	if !rg.LocConns.MatchString(l.text) {
		return exits, false
	}

	parts := strings.Split(strings.Split(l.OptimizedText(), ":")[1], ",")

	for _, part := range parts {
		words := strings.Split(strings.TrimSpace(part), " ")
		if len(words) != 2 {
			return exits, false
		}

		exits[words[0]] = words[1]
	}

	return exits, true
}

// AsItemDeclaration returns the item label, noun and adjective and true if it was found
func (l Line) AsLabelNounAdjDeclaration() (label, noun, adjetive string, ok bool) {
	if !rg.LabelNounAdjDeclaration.MatchString(l.OptimizedText()) {
		return "", "", "", false
	}

	m := rg.LabelNounAdjDeclaration.FindStringSubmatch(l.OptimizedText())
	label = m[1]
	noun = m[2]
	adjetive = m[3]

	return label, noun, adjetive, true
}

func (l Line) AsConfig() (label config.Label, value string, ok bool) {
	for _, label := range config.Labels() {
		value, ok := l.GetTextForLabel(label.String())
		if ok {
			return label, value, true
		}
	}

	return config.UnknownLabel, value, false
}
