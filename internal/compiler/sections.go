package compiler

import "strings"

type section string

const (
	sectionVars     section = "VARS"
	sectionWords    section = "WORDS"
	sectionUserMsgs section = "USER MESSAGES"
	sectionSysMsg   section = "SYSTEM MESSAGES"
	sectionObjs     section = "OBJECTS"
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
	case "OBJECTS":
		return sectionObjs
	case "LOCATIONS":
		return sectionLocs
	case "PROCESS TABLES":
		return sectionProcs
	default:
		return sectionNone
	}
}
