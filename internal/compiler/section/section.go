package section

import "strings"

type Section int

const (
	None Section = iota
	Config
	Vars
	Words
	UserMsgs
	SysMsg
	Items
	Locs
	Procs
)

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
