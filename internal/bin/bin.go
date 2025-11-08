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

func (b *BinDB) Save(path string) (int64, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	f, err := os.Create(path)
	if err != nil {
		return 0, err
	}

	gzw := gzip.NewWriter(f)

	defer func() {
		_ = gzw.Close()
		_ = f.Close()
	}()

	return b.buf.WriteTo(gzw)
}

func (b *BinDB) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}

	defer func() {
		_ = gzr.Close()
		_ = f.Close()
	}()

	b.buf.Reset()
	if _, err := io.Copy(&b.buf, gzr); err != nil {
		return err
	}

	return nil
}
