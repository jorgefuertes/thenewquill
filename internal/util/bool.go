package util

func BoolToInt(b bool) int {
	if b {
		return 1
	}

	return 0
}

func IntToBool(i int) bool {
	return i == 1
}
