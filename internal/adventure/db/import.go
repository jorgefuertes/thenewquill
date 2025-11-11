package db

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"strings"

	"github.com/fxamacker/cbor/v2"
	"github.com/jorgefuertes/thenewquill/pkg/log"
	"github.com/jorgefuertes/thenewquill/pkg/tms"
)

func (d *DB) Import(path string) error {
	d.Reset()

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Warning("error closing file %q: %s", path, err)
		}
	}()

	var headers []string
	var line string
	for {
		b := make([]byte, 1)
		if n, err := f.Read(b); err != nil {
			log.Debug("Read %d bytes", n)
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}

			return err
		}

		if string(b[0]) == "\n" {
			if line == dataSeparator {
				break
			}

			headers = append(headers, line)
			line = ""

			continue
		}

		line += string(b[0])
	}

	var remaining []byte
	remaining, err = io.ReadAll(f)
	if err != nil {
		return err
	}

	// get the key
	key := tms.GenerateKey(strings.Join(headers, "\n"))

	// decrypt
	remaining, err = tms.Decrypt(key, remaining)
	if err != nil {
		return err
	}

	// decompress
	r, err := gzip.NewReader(bytes.NewReader(remaining))
	if err != nil {
		return err
	}

	defer func() {
		if err := r.Close(); err != nil {
			log.Warning("error closing gzip reader: %s", err)
		}
	}()

	cbData, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return cbor.Unmarshal(cbData, d)
}
