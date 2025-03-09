package db_test

import (
	"bytes"
	"testing"

	"thenewquill/internal/compiler/db"
	"thenewquill/internal/compiler/section"

	"github.com/stretchr/testify/require"
)

func TestEncodeAndDecode(t *testing.T) {
	db := db.NewDB()

	headers := []string{"This is a test", "A database test"}
	regs := []struct {
		section section.Section
		label   string
		value   any
	}{
		{section.Config, "version", 2},
		{section.Vars, "test", true},
		{section.Vars, "aFloat", 1.123},
		{section.Vars, "aString", "This is a string"},
	}

	db.AddHeader(headers...)
	for _, f := range regs {
		db.AddReg(f.section, f.label, f.value)
	}

	w := new(bytes.Buffer)
	require.NotNil(t, w)

	err := db.Write(w)
	require.NoError(t, err)
	require.NotZero(t, w.Len())

	r := bytes.NewReader(w.Bytes())
	require.NotNil(t, r)

	db.Reset()
	err = db.Load(r)
	require.NoError(t, err)

	for _, h := range headers {
		require.Contains(t, db.GetHeaders(), h)
	}

	require.Len(t, db.GetRegs(), len(regs))
}
