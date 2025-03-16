package config

import "strings"

type Lang int8

const (
	Undefined = 0
	ES        = 1
	EN        = 2
)

func (l Lang) Int() int {
	return int(l)
}

func (l Lang) String() string {
	switch l {
	case ES:
		return "ES"
	case EN:
		return "EN"
	default:
		return "UNDEFINED"
	}
}

func LangFromString(s string) Lang {
	switch strings.ToUpper(s) {
	case "ES":
		return ES
	case "EN":
		return EN
	default:
		return Undefined
	}
}
