package db_test

import (
	"bytes"
	"testing"

	"thenewquill/internal/compiler/db"
	"thenewquill/internal/compiler/section"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDB(t *testing.T) {
	d := db.New()

	d.AddHeader("This is a test")
	d.AddHeader("A database test")

	d.Append(section.Vars, "aString", "This is a string")
	d.Append(section.Vars, "aInt", 10)
	d.Append(section.Vars, "aFloat", 1.123)
	d.Append(section.Vars, "aBoolTrue", true)
	d.Append(section.Vars, "aBoolFalse", false)
	d.Append(section.Vars, "anArray", []string{"A", "B", "C"})
	d.Append(section.Config, "multi", 1, 2, 3, "A", "B", "C")

	// save
	w := new(bytes.Buffer)
	require.NoError(t, d.Save(w))
	require.NotZero(t, w.Len())

	// load
	r := bytes.NewReader(w.Bytes())
	d2 := db.New()
	require.NoError(t, d2.Load(r))

	// hash
	h1, err := d.Hash()
	require.NoError(t, err)
	h2, err := d2.Hash()
	require.NoError(t, err)
	require.Equal(t, h1, h2)

	// compare headers
	for i, h := range d.Headers {
		assert.Equal(t, h, d2.Headers[i])
	}

	// compare records
	for i, r := range d.Records {
		assert.Equal(t, r, d2.Records[i])
	}
}
