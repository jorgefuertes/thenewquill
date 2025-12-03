package message

import (
	"fmt"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
)

type Message struct {
	ID      primitive.ID
	LabelID primitive.ID
	Text    string
	Plurals [2]string
}

var _ adapter.Storeable = &Message{}

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

func New(id primitive.ID, text string) Message {
	m := Message{ID: id, Text: text, Plurals: [2]string{}}

	return m
}

func (m *Message) SetID(id primitive.ID) {
	m.ID = id
}

func (m Message) GetID() primitive.ID {
	return m.ID
}

func (m *Message) SetLabelID(id primitive.ID) {
	m.LabelID = id
}

func (m Message) GetLabelID() primitive.ID {
	return m.LabelID
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
