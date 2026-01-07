package blob

import (
	"mime"
	"os"
	"path"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database/adapter"
)

type Blob struct {
	ID      uint32
	LabelID uint32 `valid:"required"`
	Mime    string `valid:"required"`
	Data    []byte `valid:"required"`
}

var _ adapter.Storeable = &Blob{}

func New() *Blob {
	b := &Blob{Data: make([]byte, 0)}

	return b
}

func (b *Blob) Load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	b.Data = make([]byte, info.Size())

	_, err = f.Read(b.Data)
	if err != nil {
		return err
	}

	b.Mime = mime.TypeByExtension(path.Ext(filename))

	return nil
}

func (b *Blob) GetKind() kind.Kind {
	return kind.Blob
}

func (b *Blob) SetID(id uint32) {
	b.ID = id
}

func (b *Blob) GetID() uint32 {
	return b.ID
}

func (b *Blob) SetLabelID(id uint32) {
	b.LabelID = id
}

func (b *Blob) GetLabelID() uint32 {
	return b.LabelID
}
