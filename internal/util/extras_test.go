package util_test

import (
	"strings"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/util"

	"github.com/stretchr/testify/assert"
)

func TestHumanBytes(t *testing.T) {
	testCases := []struct {
		in   int64
		want string
	}{
		{0, "0 B"},
		{1, "1 B"},
		{1023, "1023 B"},
		{1024, "1.0 KiB"},
		{1536, "1.5 KiB"},
		{1024 * 1024, "1.0 MiB"},
		{int64(1024) * 1024 * 1024, "1.0 GiB"},
		{int64(1024) * 1024 * 1024 * 1024, "1.0 TiB"},
	}

	for _, tc := range testCases {
		t.Run(tc.want, func(t *testing.T) {
			assert.Equal(t, tc.want, util.HumanBytes(tc.in))
		})
	}
}

func TestStringToInt(t *testing.T) {
	assert.Equal(t, 0, util.StringToInt(""))
	assert.Equal(t, 0, util.StringToInt("not a number"))
	assert.Equal(t, 42, util.StringToInt("42"))
	assert.Equal(t, -17, util.StringToInt("-17"))
}

func TestEscapeAndSplitFieldsRoundtrip(t *testing.T) {
	// EscapeField is meant to be the inverse of the decoder inside
	// SplitIntoFields for the entry that carries the @B64: prefix.
	original := []string{
		"plain text",
		"with | pipe",
		`includes "quotes" and 'apostrophes'`,
		"acentos: ñoño café",
	}

	parts := make([]string, len(original))
	for i, s := range original {
		parts[i] = util.EscapeField(s)
	}

	got := util.SplitIntoFields(strings.Join(parts, "|"))
	assert.Equal(t, original, got, "EscapeField then SplitIntoFields should be a roundtrip")
}

func TestSplitIntoFieldsLeavesPlainFieldsAlone(t *testing.T) {
	assert.Equal(t, []string{"a", "b", "c"}, util.SplitIntoFields("a|b|c"))
}

func TestSplitIntoFieldsKeepsInvalidBase64AsIs(t *testing.T) {
	// A field that starts with @B64: but is not valid base64 should be
	// returned untouched rather than producing garbage.
	assert.Equal(t, []string{"@B64:not-base64-data!!!"},
		util.SplitIntoFields("@B64:not-base64-data!!!"))
}
