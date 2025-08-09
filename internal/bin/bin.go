package bin

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"sync"
)

type BinDB struct {
	buf  bytes.Buffer
	lock sync.Mutex
}

func New() *BinDB {
	return &BinDB{
		buf:  bytes.Buffer{},
		lock: sync.Mutex{},
	}
}

func (b *BinDB) GetReader() io.Reader {
	b.lock.Lock()
	defer b.lock.Unlock()

	return bytes.NewReader(b.buf.Bytes())
}

func (b *BinDB) Save(path string) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	gzw := gzip.NewWriter(f)
	defer gzw.Close()

	_, err = b.buf.WriteTo(gzw)
	if err != nil {
		return err
	}

	return nil
}

func (b *BinDB) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	b.buf.Reset()
	_, err = io.Copy(&b.buf, gzr)
	if err != nil {
		return err
	}

	return nil
}
