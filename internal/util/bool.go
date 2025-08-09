package util

func BoolToByte(b bool) byte {
	if b {
		return 1
	}

	return 0
}

func ByteToBool(i int) bool {
	return i == 1
}
