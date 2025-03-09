package section

import (
	"slices"
	"strings"
)

type Section int

const (
	None     Section = 0
	Config   Section = 1
	Vars     Section = 2
	Words    Section = 3
	Messages Section = 4
	Items    Section = 5
	Locs     Section = 6
	Procs    Section = 7
	Chars    Section = 8
)

func Sections() []Section {
	sections := make([]Section, 0)

	for sec := range sectionNamesAndAliases() {
		sections = append(sections, sec)
	}

	return sections
}

func sectionNamesAndAliases() map[Section][]string {
	return map[Section][]string{
		None:     {"none", "unknown"},
		Config:   {"config", "cfg", "configuration"},
		Vars:     {"vars", "variables"},
		Words:    {"words", "vocabulary", "voc"},
		Messages: {"messages", "msgs"},
		Items:    {"items", "objects"},
		Locs:     {"locations", "rooms", "locs"},
		Procs:    {"process tables", "procs", "proc tables"},
		Chars:    {"characters", "chars", "char"},
	}
}

func (s Section) String() string {
	names, ok := sectionNamesAndAliases()[s]
	if !ok {
		return sectionNamesAndAliases()[None][0]
	}

	return names[0]
}

func (s Section) ToInt() int {
	return int(s)
}

func FromInt(i int) Section {
	if i < 0 || i >= len(sectionNamesAndAliases()) {
		return None
	}

	return Section(i)
}

func FromString(s string) Section {
	s = strings.ToLower(s)

	for sec, names := range sectionNamesAndAliases() {
		if slices.Contains(names, s) {
			return sec
		}
	}

	return None
}
