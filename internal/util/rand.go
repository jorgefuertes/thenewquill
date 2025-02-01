package util

import "math/rand"

func RandomString(maxLen int) string {
	chars := []byte(`abcdefghijklmnopqrstuwxyzABCDEFGHIJKLMNOPQRSTUWXYZ`)
	var out []byte
	for range maxLen {
		i := rand.Intn(len(chars))
		out = append(out, chars[i])
	}

	return string(out)
}
