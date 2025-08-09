package kind

import (
	"reflect"
	"slices"
	"strings"
)

type Kind byte

const (
	None      Kind = 0
	Label     Kind = 1
	Config    Kind = 2
	Variable  Kind = 3
	Word      Kind = 4
	Message   Kind = 5
	Item      Kind = 6
	Location  Kind = 7
	Process   Kind = 8
	Character Kind = 9
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
		Config:    {"config", "cfg", "configuration"},
		Variable:  {"var", "variable"},
		Word:      {"word", "vocabulary", "voc"},
		Message:   {"message", "msg"},
		Item:      {"item", "object"},
		Location:  {"location", "room", "loc"},
		Character: {"character", "char", "player"},
		Process:   {"process table", "proc", "proc table"},
		Label:     {"label"},
	}
}

func (s Kind) String() string {
	names, ok := kindNamesAndAliases()[s]
	if !ok {
		return kindNamesAndAliases()[None][0]
	}

	return names[0]
}

func (s Kind) Byte() byte {
	return byte(s)
}

func KindFromByte(b byte) Kind {
	if int(b) < 0 || int(b) >= len(kindNamesAndAliases()) {
		return None
	}

	return Kind(b)
}

func KindFromString(s string) Kind {
	s = strings.ToLower(s)

	for sec, names := range kindNamesAndAliases() {
		if slices.Contains(names, s) {
			return sec
		}
	}

	return None
}

func KindOf(s any) Kind {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return KindFromString(t.Name())
}
