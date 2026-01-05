package database

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database/adapter"
)

type Record struct {
	LabelID uint32
	Kind    kind.Kind
	Data    []byte
}

func (r *Record) Unmarshal(dest adapter.Storeable) error {
	return cbor.Unmarshal(r.Data, dest)
}

func (r *Record) Marshal(entity adapter.Storeable) error {
	var err error
	r.Data, err = cbor.Marshal(entity)

	return err
}
