package util

const (
	TrueSring   = "T"
	FalseString = "F"
)

func BoolToString(b bool) string {
	if b {
		return TrueSring
	}

	return FalseString
}

func StringToBool(s string) bool {
	return s == TrueSring
}
