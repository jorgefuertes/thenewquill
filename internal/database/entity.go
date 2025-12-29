package database

import (
	"fmt"
	"reflect"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

const (
	IDFieldName      = "ID"
	LabelIDFieldName = "LabelID"
)

func checkEntity(e any) (id, labelID uint32) {
	if kind.KindOf(e) == kind.None {
		panic(ErrCannotInferKind)
	}

	v := reflect.ValueOf(e)

	if v.Kind() != reflect.Ptr {
		panic(ErrEntityIsNotPointer)
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		panic(ErrEntityIsNotPointerToStruct)
	}

	idField := v.FieldByName(IDFieldName)
	if !idField.IsValid() {
		panic(ErrMissingIDField)
	}

	if idField.Kind() != reflect.Uint32 {
		panic(fmt.Errorf("'ID' field is %v, expected uint32", idField.Kind()))
	}

	if !idField.CanSet() {
		panic(ErrCannotSetIDField)
	}

	id = uint32(idField.Uint())

	labelIDField := v.FieldByName(LabelIDFieldName)
	if !labelIDField.IsValid() {
		panic(ErrMissingLabelIDField)
	}

	if labelIDField.Kind() != reflect.Uint32 {
		panic(fmt.Errorf("'LabelID' field is %v, expected uint32", labelIDField.Kind()))
	}

	labelID = uint32(labelIDField.Uint())

	return id, labelID
}

func setID(e any, id uint32) {
	_, _ = checkEntity(e)

	v := reflect.ValueOf(e).Elem()
	idField := v.FieldByName(IDFieldName)
	idField.SetUint(uint64(id))
}
