package database

import (
	"compress/gzip"
	"io"
	"os"
	"regexp"
)

type fileReader struct {
	filename string
	version  string
	fileType string
	raw      *os.File
	z        *gzip.Reader
}

func (f *fileReader) close() error {
	if err := f.z.Close(); err != nil {
		return err
	}

	return f.raw.Close()
}

func newFileReader(filename string) (*fileReader, error) {
	f := &fileReader{filename: filename}

	var err error
	f.raw, err = os.OpenFile(filename, os.O_RDONLY, 0o644)
	if err != nil {
		return f, err
	}

	f.version, f.fileType, err = f.readHeaders()
	if err != nil {
		return f, err
	}

	f.z, err = gzip.NewReader(f.raw)
	if err != nil {
		return f, err
	}

	return f, nil
}

func (f *fileReader) readLn() (string, error) {
	var reader io.Reader

	reader = f.raw
	if f.z != nil {
		reader = f.z
	}

	line := ""
	for {
		b := make([]byte, 1)
		if _, err := reader.Read(b); err != nil {
			return line, err
		}

		if b[0] == '\n' {
			break
		}

		line += string(b)
	}

	return line, nil
}

// readHeaders reads the headers of the file and returns the format version
func (f *fileReader) readHeaders() (string, string, error) {
	formatVersionRg := regexp.MustCompile(
		`^# Format version: ([\d\.]+), type: (` + databaseType + `|` + saveType + `)\s+#$`,
	)
	ver := ""
	fileType := ""

	for {
		line, err := f.readLn()
		if err != nil {
			return ver, fileType, err
		}

		if line == "---" {
			break
		}

		if formatVersionRg.MatchString(line) {
			matches := formatVersionRg.FindAllStringSubmatch(line, 1)

			ver = matches[0][1]
			fileType = matches[0][2]
		}
	}

	return ver, fileType, nil
}
