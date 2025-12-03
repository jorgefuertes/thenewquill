package adapter

import "github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"

type Storeable interface {
	GetID() primitive.ID
	SetID(id primitive.ID)
	GetLabelID() primitive.ID
	SetLabelID(id primitive.ID)
	Validate(allowNoID bool) error
}
