package section

import "strings"

type Section int

const (
	Vars Section = iota
	Words
	UserMsgs
	SysMsg
	Items
	Locs
	Procs
	None
)

func sectionNames() []string {
	return []string{
		"vars",
		"words",
		"user messages",
		"system messages",
		"items",
		"locations",
		"process tables",
		"none",
	}
}

func (s Section) String() string {
	return sectionNames()[s]
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
