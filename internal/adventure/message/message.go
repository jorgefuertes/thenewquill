package message

import (
	"fmt"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/id"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

type Message struct {
	ID      id.ID
	Text    string
	Plurals [2]string
}

var _ adapter.Storeable = Message{}

type Plural int

const (
	Zero Plural = -1
	One  Plural = 0
	Many Plural = 1
)

func PluralFromString(s string) Plural {
	switch s {
	case "one":
		return One
	case "many":
		return Many
	default:
		return Zero
	}
}

func New(text string) Message {
	m := Message{ID: id.Undefined, Text: text, Plurals: [2]string{}}

	return m
}

func (m Message) SetID(id id.ID) adapter.Storeable {
	m.ID = id

	return m
}

func (m Message) GetID() id.ID {
	return m.ID
}

func (m Message) GetKind() kind.Kind {
	return kind.Message
}

func (m *Message) SetPlural(pluralName Plural, text string) {
	if pluralName == Zero {
		m.Text = text

		return
	}

	m.Plurals[pluralName] = text
}

func (m *Message) SetPlurals(plurals [2]string) {
	m.Plurals = plurals
}

func (m Message) String() string {
	return m.Text
}

func (m Message) CountPlaceholders() int {
	return strings.Count(m.Text, "_")
}

func (m Message) IsPluralized() bool {
	return m.Plurals[One] != "" || m.Plurals[Many] != ""
}

func (m Message) Stringf(args ...any) string {
	if m.IsPluralized() && len(args) > 0 {
		return m.pluralize(args[0])
	}

	if m.CountPlaceholders() == 0 {
		return m.Text
	}

	if len(args) > m.CountPlaceholders() {
		args = args[:m.CountPlaceholders()]
	}

	for len(args) < m.CountPlaceholders() {
		args = append(args, "?")
	}

	format := strings.ReplaceAll(m.Text, "_", "%v")
	return fmt.Sprintf(format, args...)
}

func (m Message) pluralize(arg any) string {
	switch arg := arg.(type) {
	case int:
		switch arg {
		case 0:
			return m.Text
		case 1:
			return m.Plurals[One]
		default:
			return sprintf(m.Plurals[Many], arg)
		}
	case float64:
		return sprintf(m.Plurals[Many], fmt.Sprintf("%.2f", arg))
	case string:
		switch arg {
		case "0", "zero", "cero":
			return m.Text
		case "1", "one", "un", "uno", "una":
			return m.Plurals[One]
		default:
			return sprintf(m.Plurals[Many], arg)
		}
	default:
		return sprintf(m.Plurals[Many], arg)
	}
}

func sprintf(format string, args ...any) string {
	if strings.Contains(format, "_") {
		return fmt.Sprintf(strings.ReplaceAll(format, "_", "%v"), args...)
	}

	return format
}
