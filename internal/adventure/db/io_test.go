package db

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBinary(t *testing.T) {
	w := bytes.NewBuffer(nil)

	var i1 int32 = -math.MaxInt32
	err := binary.Write(w, endian, i1)
	require.NoError(t, err)

	var i12 int32
	err = binary.Read(bytes.NewReader(w.Bytes()), endian, &i12)
	require.NoError(t, err)
	require.Equal(t, i1, i12)
}

func TestWriteFields(t *testing.T) {
	fields := []any{
		byte(255),
		uint16(65535),
		int32(-math.MaxInt32),
		int64(-math.MaxInt64),
		float32(math.MaxFloat32),
		float64(math.MaxFloat64),
		"test string",
	}

	w := bytes.NewBuffer(nil)

	for _, f := range fields {
		write(w, f)
	}

	r := bytes.NewReader(w.Bytes())

	for _, field := range fields {
		switch orig := field.(type) {
		case byte:
			var b byte
			read(r, &b)
			require.Equal(t, orig, b)
		case uint16:
			var u uint16
			read(r, &u)
			require.Equal(t, orig, u)
		case int32:
			var i int32
			read(r, &i)
			require.Equal(t, orig, i)
		case int64:
			var i int64
			read(r, &i)
			require.Equal(t, orig, i)
		case float32:
			var f float32
			read(r, &f)
			require.Equal(t, orig, f)
		case float64:
			var f float64
			read(r, &f)
			require.Equal(t, orig, f)
		case string:
			var s string
			read(r, &s)
			require.Equal(t, orig, s)
		default:
			t.Error("unknown type")
		}
	}
}
