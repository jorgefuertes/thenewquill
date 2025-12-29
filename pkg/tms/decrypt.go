package tms

import (
	"crypto/aes"
	"crypto/cipher"
)

func Decrypt(key, encrypted []byte) ([]byte, error) {
	if !IsValidKey(key) {
		return nil, ErrInvalidKey
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encrypted) < nonceSize {
		return nil, ErrBufferTooShort
	}

	nonce, encrypted := encrypted[:nonceSize], encrypted[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
