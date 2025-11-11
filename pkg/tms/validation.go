package tms

func IsValidKey(key []byte) bool {
	return len(key) == 32
}