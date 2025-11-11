package tms

import (
	"crypto/sha256"
	"regexp"
)

func GenerateKey(seed string) []byte {
	rg := regexp.MustCompile(`(#|\s)+`)
	seed = rg.ReplaceAllString(seed, "")

	hash := sha256.Sum256([]byte(seed))

	return hash[:]
}