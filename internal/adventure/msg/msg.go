package msg

import (
	"fmt"
	"regexp"
	"strings"
)

type Msg struct {
	Label   string
	Text    string
	Plurals [3]string
}

var (
	PluralNames   = []string{"zero", "one", "many"}
	pluralLabelRg = regexp.MustCompile(`^([\d\p{L}\-_]+)\.(zero|one|many)$`)
)

func New(label, text string) *Msg {
	m := &Msg{Label: label, Text: text}

	if !pluralLabelRg.MatchString(label) {
		return m
	}

	matches := pluralLabelRg.FindStringSubmatch(label)
	m.Text = ""
	m.Label = matches[1]

	for i, p := range PluralNames {
		if p == matches[2] {
			m.Plurals[i] = text
		}
	}

	return m
}

func (m *Msg) GetLabel() string {
	return m.Label
}

func (m *Msg) SetPlurals(texts [3]string) {
	m.Plurals = texts
}

func (m Msg) String() string {
	return m.Text
}

func (m Msg) IsPluralized() bool {
	for _, p := range m.Plurals {
		if p != "" {
			return true
		}
	}

	return false
}

func (m Msg) Stringf(args ...any) string {
	if m.IsPluralized() && len(args) > 0 {
		return pluralize(m.Plurals, args[0])
	}

	format := strings.ReplaceAll(m.Text, "_", "%v")
	return fmt.Sprintf(format, args...)
}

func pluralize(texts [3]string, arg any) string {
	switch arg := arg.(type) {
	case int:
		switch arg {
		case 0:
			return sprintf(texts[0], 0)
		case 1:
			return sprintf(texts[1], 1)
		default:
			return sprintf(texts[2], arg)
		}
	case float64:
		switch arg {
		case 0:
			return sprintf(texts[0], 0)
		case 1:
			return sprintf(texts[1], 1)
		default:
			return sprintf(texts[2], fmt.Sprintf("%.2f", arg))
		}
	case string:
		switch arg {
		case "0", "zero", "cero":
			return sprintf(texts[0], 0)
		case "1", "one", "un", "uno", "una":
			return sprintf(texts[1], 1)
		default:
			return sprintf(texts[2], arg)
		}
	default:
		return sprintf(texts[2], arg)
	}
}

func sprintf(format string, args ...any) string {
	if strings.Contains(format, "_") {
		return fmt.Sprintf(strings.ReplaceAll(format, "_", "%v"), args...)
	}

	return format
}
