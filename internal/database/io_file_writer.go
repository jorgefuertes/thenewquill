package database

import (
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

const zBegin = "---"

type fileWriter struct {
	filename  string
	raw       *os.File
	z         *gzip.Writer
	byteCount int64
}

// close closes the file writer and returns the number of bytes written and the final file size
func (f *fileWriter) close() (int64, int64, error) {
	if err := f.z.Flush(); err != nil {
		return f.byteCount, 0, err
	}

	if err := f.z.Close(); err != nil {
		return f.byteCount, 0, err
	}

	if err := f.raw.Sync(); err != nil {
		return f.byteCount, 0, err
	}

	info, err := f.raw.Stat()
	if err != nil {
		return f.byteCount, 0, err
	}

	if err := f.raw.Close(); err != nil {
		return f.byteCount, 0, err
	}

	return f.byteCount, info.Size(), nil
}

func newFileWriter(filename string, fileType, version byte, headers ...string) (*fileWriter, error) {
	f := &fileWriter{filename: filename}

	var err error
	f.raw, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}

	if err := f.raw.Truncate(0); err != nil {
		return nil, err
	}

	if len(headers) > 0 {
		f.writeHumanHeaders(headers)
	}

	f.z, err = gzip.NewWriterLevel(f.raw, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	f.write([]byte{fileType, version})

	return f, nil
}

func (f *fileWriter) write(data any) {
	f.byteCount += int64(binary.Size(data))

	if f.z != nil {
		if err := binary.Write(f.z, endian, data); err != nil {
			panic(fmt.Sprintf("[fz] cannot write data, %s: %+v", err, data))
		}

		return
	}

	if err := binary.Write(f.raw, endian, data); err != nil {
		panic(fmt.Sprintf("[fr] cannot write data, %s: %+v", err, data))
	}
}

func (f *fileWriter) writeString(format string, args ...any) {
	s := format
	if len(args) > 0 {
		s = fmt.Sprintf(format, args...)
	}

	f.write([]byte(s))
}

func (f *fileWriter) writeStringLn(format string, args ...any) {
	f.writeString(format+"\n", args...)
}

func (f *fileWriter) writeHumanHeaders(headers []string) {
	max := 0
	for _, h := range headers {
		if len(h) > max {
			max = len(h)
		}
	}
	max += 4

	f.writeStringLn(strings.Repeat("#", max))

	for _, h := range headers {
		f.writeStringLn("# %s%s#", h, strings.Repeat(" ", max-len(h)-3))
	}

	f.writeStringLn(strings.Repeat("#", max))
	f.writeStringLn("")
	f.writeStringLn(zBegin)
}
