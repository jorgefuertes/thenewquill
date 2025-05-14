package db_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/stretchr/testify/assert"
)

func TestSubKinds(t *testing.T) {
	subKinds := db.SubKinds()
	assert.Len(t, subKinds, 8)

	for _, subKind := range subKinds {
		assert.NotEmpty(t, subKind.String())
		assert.NotEmpty(t, subKind.Byte())
	}

	name := db.SubKind(99).String()
	assert.Equal(t, db.NoSubKind.String(), name)

	subKind := db.SubKindFromByte(0)
	assert.Equal(t, db.NoSubKind, subKind)

	subKind = db.SubKindFromByte(1)
	assert.Equal(t, db.VerbSubKind, subKind)

	subKind = db.SubKindFromByte(255)
	assert.Equal(t, db.NoSubKind, subKind)

	subKind = db.SubKindFromString("adjective")
	assert.Equal(t, db.AdjectiveSubKind, subKind)

	subKind = db.SubKindFromString("invalid")
	assert.Equal(t, db.NoSubKind, subKind)
}
