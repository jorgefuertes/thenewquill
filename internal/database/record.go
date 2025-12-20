package database

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

type Record struct {
	LabelID uint32
	Kind    kind.Kind
	Data    []byte
}

func (r *Record) Unmarshal(dest any) error {
	_, _, err := checkEntity(dest)
	if err != nil {
		return err
	}

	return cbor.Unmarshal(r.Data, dest)
}

func (r *Record) Marshal(entity any) error {
	_, _, err := checkEntity(entity)
	if err != nil {
		return err
	}

	r.Data, err = cbor.Marshal(entity)

	return err
}
