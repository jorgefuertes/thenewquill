package line

import (
	"strconv"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adventure/config"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/adventure/message"
	"github.com/jorgefuertes/thenewquill/internal/compiler/rg"
)

// AsInclude returns the include filename string and true if it was found
func (l Line) AsInclude() (string, bool) {
	if !rg.Include.MatchString(l.OptimizedText()) {
		return "", false
	}

	return rg.Include.FindStringSubmatch(l.OptimizedText())[1], true
}

// AsSection returns the section and true if it was found
func (l Line) AsSection() (kind.Kind, bool) {
	if !rg.Section.MatchString(l.OptimizedText()) {
		return kind.None, false
	}

	return kind.KindFromString(rg.Section.FindStringSubmatch(l.OptimizedText())[1]), true
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

// AsWord returns the word name, synonyms and true if it was found
func (l Line) AsWord() (string, []string, bool) {
	o := l.OptimizedText()

	if !rg.Word.MatchString(o) {
		return "", nil, false
	}

	parts := strings.Split(o, ":")
	if len(parts) != 2 {
		return "", nil, false
	}

	labelName := parts[0]

	syns := strings.Split(parts[1], ",")
	for i, syn := range syns {
		syns[i] = strings.TrimSpace(syn)
		if syns[i] == "" {
			syns = append(syns[:i], syns[i+1:]...)
			continue
		}
	}

	return labelName, syns, true
}

// AsMsg returns the label name, text, plural name and result
func (l Line) AsMsg() (string, string, message.Plural, bool) {
	if !rg.Msg.MatchString(l.text) {
		return "", "", message.Zero, false
	}

	// is a plural?
	if rg.MsgPlural.MatchString(l.text) {
		parts := rg.MsgPlural.FindStringSubmatch(l.text)
		if len(parts) != 4 {
			return "", "", message.Zero, false
		}

		return parts[2], parts[3], message.PluralFromString(parts[1]), true
	}

	parts := rg.Msg.FindStringSubmatch(l.text)
	if len(parts) != 3 {
		return "", "", message.Zero, false
	}

	return parts[1], parts[2], message.Zero, true
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
func (l Line) AsLabelNounAdjDeclaration() (labelName, nounName, adjetiveName string, ok bool) {
	if !rg.LabelNounAdjDeclaration.MatchString(l.OptimizedText()) {
		return "", "", "", false
	}

	m := rg.LabelNounAdjDeclaration.FindStringSubmatch(l.OptimizedText())
	labelName = m[1]
	nounName = m[2]
	adjetiveName = m[3]

	return labelName, nounName, adjetiveName, true
}

// AsConfig returns the config field name and value and true if it was found
func (l Line) AsConfig() (primitive.Label, string, bool) {
	for _, label := range config.AllowedFieldLabels() {
		v, ok := l.GetTextForLabel(label)

		if ok {
			return label, v, true
		}
	}

	return "", "", false
}
