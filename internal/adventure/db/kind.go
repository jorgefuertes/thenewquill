package db

import (
	"slices"
	"strings"
)

type Kind byte

const (
	None       Kind = 0
	Config     Kind = 1
	Variables  Kind = 2
	Words      Kind = 3
	Messages   Kind = 4
	Items      Kind = 5
	Locations  Kind = 6
	Processes  Kind = 7
	Characters Kind = 8
	Labels     Kind = 255
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
		None:       {"none", "unknown"},
		Config:     {"config", "cfg", "configuration"},
		Variables:  {"vars", "variables"},
		Words:      {"words", "vocabulary", "voc"},
		Messages:   {"messages", "msgs"},
		Items:      {"items", "objects"},
		Locations:  {"locations", "rooms", "locs"},
		Characters: {"characters", "chars", "players"},
		Processes:  {"process tables", "procs", "proc tables"},
		Labels:     {"labels", "labels"},
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
