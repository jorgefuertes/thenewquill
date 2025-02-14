package compiler

import "strings"

type section string

const (
	sectionVars     section = "VARS"
	sectionWords    section = "WORDS"
	sectionUserMsgs section = "USER MESSAGES"
	sectionSysMsg   section = "SYSTEM MESSAGES"
	sectionItems    section = "ITEMS"
	sectionLocs     section = "LOCATIONS"
	sectionProcs    section = "PROCESS TABLES"
	sectionNone     section = "NONE"
)

func (s section) String() string {
	return string(s)
}

func sectionFromString(s string) section {
	s = strings.ToUpper(s)

	switch s {
	case "VARS":
		return sectionVars
	case "WORDS":
		return sectionWords
	case "SYSTEM MESSAGES":
		return sectionSysMsg
	case "USER MESSAGES":
		return sectionUserMsgs
	case "ITEMS", "OBJECTS":
		return sectionItems
	case "LOCATIONS":
		return sectionLocs
	case "PROCESS TABLES":
		return sectionProcs
	default:
		return sectionNone
	}
}

func (s section) singleString() string {
	switch s {
	case sectionVars:
		return "var"
	case sectionWords:
		return "word"
	case sectionLocs:
		return "location"
	case sectionItems:
		return "item"
	case sectionProcs:
		return "process"
	case sectionSysMsg:
		return "sysmsg"
	case sectionUserMsgs:
		return "usermsg"
	default:
		return "unknown"
	}
}
