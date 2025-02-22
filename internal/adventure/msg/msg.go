package msg

import (
	"fmt"
	"regexp"
	"strings"
)

type Msg struct {
	Type   MsgType
	Label  string
	Text   string
	plural [3]string
}

var (
	plurals       = []string{"zero", "one", "more"}
	pluralLabelRg = regexp.MustCompile(`^([\d\p{L}\-_]+)\.(zero|one|more)$`)
)

const pluralized = "#plural#"

func New(t MsgType, label, text string) *Msg {
	m := &Msg{Type: t, Label: label, Text: text}

	if !pluralLabelRg.MatchString(label) {
		return m
	}

	m.Text = pluralized
	matches := pluralLabelRg.FindStringSubmatch(label)
	m.Label = matches[1]

	for i, p := range plurals {
		if p == matches[2] {
			m.plural[i] = text
		}
	}

	return m
}

func (m *Msg) SetPluralTexts(texts [3]string) {
	m.plural = texts
}

func (m Msg) String() string {
	return m.Text
}

func (m Msg) Stringf(args ...any) string {
	if m.Text == pluralized && len(args) > 0 {
		return pluralize(m.plural, args[0])
	}

	format := strings.Replace(m.Text, "_", "%v", -1)
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
		switch {
		case arg == 0:
			return sprintf(texts[0], 0)
		case arg == 1:
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
		return fmt.Sprintf(strings.Replace(format, "_", "%v", -1), args...)
	}

	return format
}
