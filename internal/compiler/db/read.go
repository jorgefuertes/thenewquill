package db

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"errors"
	"io"
	"strconv"
	"strings"

	"thenewquill/internal/compiler/section"
)

const headerLinesLimit = 500

func (db *DB) Load(r io.Reader) error {
	var err error

	db.Reset()
	br := bufio.NewReader(r)

	db.headers, err = readHeaders(br)
	if err != nil {
		return err
	}

	crunchedBuf := make([]byte, br.Buffered()-sha256.Size)
	if _, err := io.ReadFull(br, crunchedBuf); err != nil {
		return err
	}

	db.regs, err = loadRegisters(crunchedBuf)
	if err != nil {
		return err
	}

	hash := make([]byte, sha256.Size)
	_, err = br.Read(hash)
	if err != nil {
		return errors.Join(err, ErrHashMismatch)
	}

	if !bytes.Equal(hash, db.Hash()) {
		return ErrHashMismatch
	}

	return nil
}

func loadRegisters(buf []byte) ([]Register, error) {
	regs := make([]Register, 0)

	zr, err := zlib.NewReader(bytes.NewReader(buf))
	if err != nil {
		return regs, err
	}

	for {
		rawReg := make([]byte, 0)
		for {
			b := make([]byte, 1)
			_, err := zr.Read(b)
			if err != nil {
				if err == io.EOF {
					break
				}

				return regs, errors.Join(err, ErrUnexpectedEOF)
			}

			if b[0] == regSep {
				break
			}

			rawReg = append(rawReg, b[0])
		}

		if len(rawReg) == 0 {
			break
		}

		fields := strings.Split(string(rawReg), string(fieldSep))
		if len(fields) < 2 {
			return regs, ErrShortRegister
		}

		secInt, err := strconv.Atoi(fields[0])
		if err != nil {
			return regs, errors.Join(err, ErrInvalidSection)
		}

		regs = append(regs, Register{Section: section.FromInt(secInt), Fields: fields[1:]})
	}

	return regs, nil
}

func readHeaders(r io.Reader) ([]string, error) {
	headers := make([]string, 0)

	currentLineBytes := make([]byte, 0)
	for {
		b := make([]byte, 1)
		_, err := r.Read(b)
		if err != nil {
			if err == io.EOF {
				return headers, ErrUnexpectedEOF
			}

			return headers, err
		}

		if b[0] == beginRecords {
			break
		}

		if b[0] == '\n' {
			if len(headers) >= headerLinesLimit {
				return headers, ErrHeaderLimitReached
			}

			headers = append(headers, string(currentLineBytes))
			currentLineBytes = make([]byte, 0)

			continue
		}

		currentLineBytes = append(currentLineBytes, b[0])
	}

	return headers, nil
}
