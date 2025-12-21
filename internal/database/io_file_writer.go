package database

import (
	"compress/gzip"
	"fmt"
	"os"
	"strings"
)

type fileWriter struct {
	filename  string
	raw       *os.File
	z         *gzip.Writer
	byteCount int
}

func (f *fileWriter) close() (int, int, error) {
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

	return f.byteCount, int(info.Size()), nil
}

func createFile(filename string, headers ...string) (*fileWriter, error) {
	f := &fileWriter{filename: filename}

	var err error
	f.raw, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return f, err
	}

	if err := f.raw.Truncate(0); err != nil {
		return f, err
	}

	if len(headers) > 0 {
		if err := f.writeHeaders(headers); err != nil {
			return f, err
		}
	}

	f.z, err = gzip.NewWriterLevel(f.raw, gzip.BestCompression)
	if err != nil {
		return f, err
	}

	return f, nil
}

func (f *fileWriter) write(format string, args ...any) error {
	s := format
	if len(args) > 0 {
		s = fmt.Sprintf(format, args...)
	}

	if f.z != nil {
		n, err := f.z.Write([]byte(s))
		f.byteCount += n

		return err
	}

	n, err := f.raw.Write([]byte(s))
	f.byteCount += n

	return err
}

func (f *fileWriter) writeLn(format string, args ...any) error {
	format += "\n"

	return f.write(format, args...)
}

func (f *fileWriter) writeHeaders(headers []string) error {
	max := 0
	for _, h := range headers {
		if len(h) > max {
			max = len(h)
		}
	}
	max += 4

	if err := f.writeLn(strings.Repeat("#", max)); err != nil {
		return err
	}

	for _, h := range headers {
		if err := f.writeLn("# %s%s#", h, strings.Repeat(" ", max-len(h)-3)); err != nil {
			return err
		}
	}

	if err := f.writeLn(strings.Repeat("#", max)); err != nil {
		return err
	}

	return f.writeLn("---")
}
