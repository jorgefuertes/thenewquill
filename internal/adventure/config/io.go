package config

import (
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/id"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/util"
)

var _ adapter.Exportable = Param{}

const (
	KindFieldIndex = iota
	IDFieldIndex
	ValueFieldIndex
)

func (p Param) Export() string {
	return fmt.Sprintf("%d|%d|%s\n", p.GetKind().Int(), p.ID, util.EscapeField(p.V))
}

func Import(s string) (Param, error) {
	fields := util.SplitIntoFields(s)

	if !kind.Param.Is(fields[KindFieldIndex]) {
		return Param{}, fmt.Errorf("cannot import param %q: invalid kind %q", s, fields[KindFieldIndex])
	}

	return Param{ID: id.FromString(fields[IDFieldIndex]), V: fields[ValueFieldIndex]}, nil
}
