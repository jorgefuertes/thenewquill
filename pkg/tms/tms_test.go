package tms_test

import (
	"bytes"
	"testing"

	"github.com/jorgefuertes/thenewquill/pkg/tms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateKey(t *testing.T) {
	t.Run("generates 32-byte key", func(t *testing.T) {
		key := tms.GenerateKey("test seed")
		assert.Len(t, key, 32)
	})

	t.Run("same seed produces same key", func(t *testing.T) {
		key1 := tms.GenerateKey("my secret seed")
		key2 := tms.GenerateKey("my secret seed")
		assert.Equal(t, key1, key2)
	})

	t.Run("different seeds produce different keys", func(t *testing.T) {
		key1 := tms.GenerateKey("seed one")
		key2 := tms.GenerateKey("seed two")
		assert.NotEqual(t, key1, key2)
	})

	t.Run("ignores whitespace in seed", func(t *testing.T) {
		key1 := tms.GenerateKey("test")
		key2 := tms.GenerateKey("t e s t")
		key3 := tms.GenerateKey("t\te\ns\tt")
		assert.Equal(t, key1, key2)
		assert.Equal(t, key1, key3)
	})

	t.Run("ignores hash symbols in seed", func(t *testing.T) {
		key1 := tms.GenerateKey("test")
		key2 := tms.GenerateKey("t#e#s#t")
		key3 := tms.GenerateKey("###test###")
		assert.Equal(t, key1, key2)
		assert.Equal(t, key1, key3)
	})

	t.Run("empty seed produces valid key", func(t *testing.T) {
		key := tms.GenerateKey("")
		assert.Len(t, key, 32)
	})

	t.Run("generated key is valid", func(t *testing.T) {
		key := tms.GenerateKey("any seed")
		assert.True(t, tms.IsValidKey(key))
	})
}

func TestIsValidKey(t *testing.T) {
	t.Run("32-byte key is valid", func(t *testing.T) {
		key := make([]byte, 32)
		assert.True(t, tms.IsValidKey(key))
	})

	t.Run("shorter key is invalid", func(t *testing.T) {
		key := make([]byte, 16)
		assert.False(t, tms.IsValidKey(key))
	})

	t.Run("longer key is invalid", func(t *testing.T) {
		key := make([]byte, 64)
		assert.False(t, tms.IsValidKey(key))
	})

	t.Run("nil key is invalid", func(t *testing.T) {
		assert.False(t, tms.IsValidKey(nil))
	})

	t.Run("empty key is invalid", func(t *testing.T) {
		assert.False(t, tms.IsValidKey([]byte{}))
	})
}

func TestEncrypt(t *testing.T) {
	validKey := tms.GenerateKey("test key")

	t.Run("encrypts data successfully", func(t *testing.T) {
		plain := []byte("hello world")
		encrypted, err := tms.Encrypt(validKey, plain)

		require.NoError(t, err)
		assert.NotNil(t, encrypted)
		assert.NotEqual(t, plain, encrypted)
		assert.Greater(t, len(encrypted), len(plain)) // encrypted includes nonce
	})

	t.Run("returns error with invalid key", func(t *testing.T) {
		invalidKey := make([]byte, 16)
		plain := []byte("hello")

		encrypted, err := tms.Encrypt(invalidKey, plain)

		assert.ErrorIs(t, err, tms.ErrInvalidKey)
		assert.Nil(t, encrypted)
	})

	t.Run("returns error with nil key", func(t *testing.T) {
		plain := []byte("hello")

		encrypted, err := tms.Encrypt(nil, plain)

		assert.ErrorIs(t, err, tms.ErrInvalidKey)
		assert.Nil(t, encrypted)
	})

	t.Run("encrypts empty data", func(t *testing.T) {
		encrypted, err := tms.Encrypt(validKey, []byte{})

		require.NoError(t, err)
		assert.NotNil(t, encrypted)
	})

	t.Run("produces different ciphertext each time", func(t *testing.T) {
		plain := []byte("same message")

		encrypted1, err1 := tms.Encrypt(validKey, plain)
		encrypted2, err2 := tms.Encrypt(validKey, plain)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, encrypted1, encrypted2) // due to random nonce
	})
}

func TestDecrypt(t *testing.T) {
	validKey := tms.GenerateKey("test key")

	t.Run("decrypts data successfully", func(t *testing.T) {
		original := []byte("secret message")
		encrypted, err := tms.Encrypt(validKey, original)
		require.NoError(t, err)

		decrypted, err := tms.Decrypt(validKey, encrypted)

		require.NoError(t, err)
		assert.Equal(t, original, decrypted)
	})

	t.Run("returns error with invalid key", func(t *testing.T) {
		invalidKey := make([]byte, 16)
		encrypted := []byte("some encrypted data")

		decrypted, err := tms.Decrypt(invalidKey, encrypted)

		assert.ErrorIs(t, err, tms.ErrInvalidKey)
		assert.Nil(t, decrypted)
	})

	t.Run("returns error with nil key", func(t *testing.T) {
		encrypted := []byte("some encrypted data")

		decrypted, err := tms.Decrypt(nil, encrypted)

		assert.ErrorIs(t, err, tms.ErrInvalidKey)
		assert.Nil(t, decrypted)
	})

	t.Run("returns error with wrong key", func(t *testing.T) {
		wrongKey := tms.GenerateKey("wrong key")
		original := []byte("secret")

		encrypted, err := tms.Encrypt(validKey, original)
		require.NoError(t, err)

		decrypted, err := tms.Decrypt(wrongKey, encrypted)

		assert.Error(t, err)
		assert.Nil(t, decrypted)
	})

	t.Run("returns error with buffer too short", func(t *testing.T) {
		shortBuffer := []byte{1, 2, 3}

		decrypted, err := tms.Decrypt(validKey, shortBuffer)

		assert.ErrorIs(t, err, tms.ErrBufferTooShort)
		assert.Nil(t, decrypted)
	})

	t.Run("returns error with tampered ciphertext", func(t *testing.T) {
		original := []byte("secret message")
		encrypted, err := tms.Encrypt(validKey, original)
		require.NoError(t, err)

		// Tamper with the ciphertext
		encrypted[len(encrypted)-1] ^= 0xFF

		decrypted, err := tms.Decrypt(validKey, encrypted)

		assert.Error(t, err)
		assert.Nil(t, decrypted)
	})

	t.Run("decrypts empty data", func(t *testing.T) {
		encrypted, err := tms.Encrypt(validKey, []byte{})
		require.NoError(t, err)

		decrypted, err := tms.Decrypt(validKey, encrypted)

		require.NoError(t, err)
		assert.Empty(t, decrypted)
	})
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	key := tms.GenerateKey("roundtrip test key")

	testCases := []struct {
		name string
		data []byte
	}{
		{"empty", []byte{}},
		{"single byte", []byte{0x42}},
		{"short string", []byte("hello")},
		{"long string", []byte("This is a longer message that should still encrypt and decrypt correctly.")},
		{"binary data", []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}},
		{"unicode", []byte("こんにちは世界 🌍")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := tms.Encrypt(key, tc.data)
			require.NoError(t, err)

			decrypted, err := tms.Decrypt(key, encrypted)
			require.NoError(t, err)

			assert.True(t, bytes.Equal(tc.data, decrypted))
		})
	}
}

func TestErrors(t *testing.T) {
	t.Run("ErrBufferTooShort has message", func(t *testing.T) {
		assert.Contains(t, tms.ErrBufferTooShort.Error(), "buffer")
	})

	t.Run("ErrInvalidKey has message", func(t *testing.T) {
		assert.Contains(t, tms.ErrInvalidKey.Error(), "key")
	})

	t.Run("errors are distinct", func(t *testing.T) {
		assert.NotEqual(t, tms.ErrBufferTooShort, tms.ErrInvalidKey)
	})
}

func TestLargeData(t *testing.T) {
	key := tms.GenerateKey("large data test")

	// 1MB of data
	largeData := make([]byte, 1024*1024)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	encrypted, err := tms.Encrypt(key, largeData)
	require.NoError(t, err)

	decrypted, err := tms.Decrypt(key, encrypted)
	require.NoError(t, err)

	assert.True(t, bytes.Equal(largeData, decrypted))
}
