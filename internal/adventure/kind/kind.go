package kind

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

type Kind byte

const (
	None Kind = iota
	Label
	Param
	Variable
	Word
	Message
	Item
	Location
	Character
	Process
	Test
)

func Kinds() []Kind {
	kinds := make([]Kind, 0)

	for sec := range kindNamesAndAliases() {
		kinds = append(kinds, sec)
	}

	return kinds
}

func kindNamesAndAliases() map[Kind][]string {
	return map[Kind][]string{
		None:      {"none", "unknown"},
		Label:     {"label", "labels"},
		Param:     {"config", "cfg", "configuration", "param", "params"},
		Variable:  {"var", "variable", "vars", "variables"},
		Word:      {"word", "vocabulary", "voc", "words"},
		Message:   {"message", "msg", "messages"},
		Item:      {"item", "object", "items", "objects"},
		Location:  {"location", "room", "loc", "locations", "rooms"},
		Character: {"character", "char", "player", "players", "characters", "chars"},
		Process:   {"process table", "proc", "proc table", "process tables", "proc tables", "procs"},
		Test:      {"test", "testitem"},
	}
}

func (s Kind) String() string {
	names, ok := kindNamesAndAliases()[s]
	if !ok {
		return kindNamesAndAliases()[None][0]
	}

	return names[0]
}

func (s Kind) TitleString() string {
	titles := map[Kind]string{
		None:      "NONE",
		Label:     "LABELS",
		Param:     "CONFIG",
		Variable:  "VARIABLES",
		Word:      "WORDS",
		Message:   "MESSAGES",
		Item:      "ITEMS",
		Location:  "LOCATIONS",
		Character: "CHARACTERS",
		Process:   "PROCESS TABLES",
		Test:      "TEST",
	}

	title, ok := titles[s]
	if !ok {
		return titles[None]
	}

	return title
}

func (s Kind) Byte() byte {
	return byte(s)
}

func (s Kind) Int() uint8 {
	return uint8(s)
}

func KindFromByte(b byte) Kind {
	if int(b) < 0 || int(b) >= len(kindNamesAndAliases()) {
		return None
	}

	return Kind(b)
}

func KindFromString(s string) Kind {
	if s == "" {
		return None
	}

	s = strings.ToLower(s)

	for sec, names := range kindNamesAndAliases() {
		if slices.Contains(names, s) {
			return sec
		}
	}

	return None
}

func KindOf(s any) Kind {
	if s == nil {
		return None
	}

	t := reflect.TypeOf(s)

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return KindFromString(t.Name())
}

func (k Kind) Is(s string) bool {
	if k.String() == s {
		return true
	}

	return fmt.Sprint(k.Int()) == s
}
