package section

import (
	"slices"
	"strings"
)

type Section int

const (
	None    Section = 0
	Config  Section = 1
	Vars    Section = 2
	Words   Section = 3
	UserMsg Section = 4
	SysMsg  Section = 5
	Items   Section = 6
	Locs    Section = 7
	Procs   Section = 8
	Chars   Section = 9
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
		None:    {"none", "unknown"},
		Config:  {"config", "cfg", "configuration"},
		Vars:    {"vars", "variables"},
		Words:   {"words", "vocabulary", "voc"},
		UserMsg: {"user messages", "user msgs", "usermsgs"},
		SysMsg:  {"system messages", "sys msgs", "system msgs", "sysmsgs"},
		Items:   {"items", "objects"},
		Locs:    {"locations", "rooms", "locs"},
		Procs:   {"process tables", "procs", "proc tables"},
		Chars:   {"characters", "chars", "char"},
	}
}

func (s Section) String() string {
	if s < 0 || s >= Section(len(sectionNamesAndAliases())) {
		return sectionNamesAndAliases()[None][0]
	}

	return sectionNamesAndAliases()[s][0]
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
