package section

import "strings"

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
	return []Section{
		None,
		Config,
		Vars,
		Words,
		UserMsg,
		SysMsg,
		Items,
		Locs,
		Procs,
		Chars,
	}
}

func sectionNames() []string {
	return []string{
		"none",
		"config",
		"vars",
		"words",
		"user messages",
		"system messages",
		"items",
		"locations",
		"process tables",
		"characters",
	}
}

func (s Section) String() string {
	if s < 0 || s >= Section(len(sectionNames())) {
		return sectionNames()[None]
	}

	return sectionNames()[s]
}

func FromInt(i int) Section {
	if i < 0 || i >= len(sectionNames()) {
		return None
	}

	return Section(i)
}

func FromString(s string) Section {
	s = strings.ToLower(s)

	for i, name := range sectionNames() {
		if name == s {
			return Section(i)
		}
	}

	return None
}
