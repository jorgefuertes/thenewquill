package db

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

var endian = binary.BigEndian

const (
	databaseSep byte = 0x1c
	recordSep   byte = 0x1e
	strBegin    byte = 0x1f
	strEnd      byte = 0x20
)

func (d *DB) Export(w io.Writer) error {
	// database start
	write(w, databaseSep)

	// labels
	for _, l := range d.Labels {
		write(w, Labels.Byte())
		write(w, l.ID.UInt32())
		write(w, l.Name)
		write(w, recordSep)
	}

	// database end
	write(w, databaseSep)

	return nil
}

func (d *DB) Import(r io.Reader) error {
	for {
		b := make([]byte, 1)
		if _, err := r.Read(b); err != nil {
			if err == io.EOF {
				break
			}

			return err
		}
	}

	return nil
}

func write(w io.Writer, f any) {
	if err := binary.Write(w, endian, f); err != nil {
		errorAndExit(err)
	}
}

func read(r io.Reader, f any) {
	if err := binary.Read(r, endian, f); err != nil {
		errorAndExit(err)
	}
}

func errorAndExit(err error) {
	fmt.Println(err)

	os.Exit(1)
}
