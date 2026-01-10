package database

import (
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"os"
)

type fileReader struct {
	filename string
	fileType byte
	version  byte
	raw      *os.File
	z        *gzip.Reader
}

func (f *fileReader) close() error {
	if err := f.z.Close(); err != nil {
		return err
	}

	return f.raw.Close()
}

// newFileReader opens the file and prepares it for reading
// reads the header fileType and version
func newFileReader(filename string) (*fileReader, error) {
	f := &fileReader{filename: filename}

	var err error
	f.raw, err = os.OpenFile(filename, os.O_RDONLY, 0o644)
	if err != nil {
		return f, err
	}

	// skip human headers
	for {
		line, err := f.readLine()
		if err != nil {
			return nil, fmt.Errorf("cannot find the zBegin mark: %s", err)
		}

		if line == zBegin {
			break
		}
	}

	f.z, err = gzip.NewReader(f.raw)
	if err != nil {
		return nil, err
	}

	var fileType byte
	if err := binary.Read(f.z, endian, &fileType); err != nil {
		return nil, err
	}

	var version byte
	if err := binary.Read(f.z, endian, &version); err != nil {
		return nil, err
	}

	f.fileType = fileType
	f.version = version

	return f, nil
}

func (f *fileReader) readLine() (string, error) {
	var line []byte
	buf := make([]byte, 1)

	for {
		var err error
		if f.z != nil {
			_, err = f.z.Read(buf)
		} else {
			_, err = f.raw.Read(buf)
		}

		if err != nil {
			return "", err
		}

		if buf[0] == '\n' {
			break
		}

		line = append(line, buf[0])
	}

	return string(line), nil
}
