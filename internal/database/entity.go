package database

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

const (
	IDFieldName      = "ID"
	LabelIDFieldName = "LabelID"
)

func checkDst(dst any) error {
	_, _, err := checkEntity(dst)

	return err
}

func checkEntity(e any) (id, labelID uint32, err error) {
	if kind.KindOf(e) == kind.None {
		return id, labelID, ErrCannotInferKind
	}

	v := reflect.ValueOf(e)

	if v.Kind() != reflect.Ptr {
		return id, labelID, errors.Join(ErrEntityIsNotPointer, errors.New("expected pointer, got "+v.Kind().String()))
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return id, labelID, errors.Join(
			ErrEntityIsNotPointerToStruct,
			errors.New("expected pointer to struct, got "+v.Kind().String()),
		)
	}

	idField := v.FieldByName(IDFieldName)
	if !idField.IsValid() {
		return id, labelID, ErrMissingIDField
	}

	if idField.Kind() != reflect.Uint32 {
		return id, labelID, fmt.Errorf("'ID' field is %v, expected uint32", idField.Kind())
	}

	if !idField.CanSet() {
		return id, labelID, ErrCannotSetIDField
	}

	id = uint32(idField.Uint())

	labelIDField := v.FieldByName(LabelIDFieldName)
	if !labelIDField.IsValid() {
		return id, labelID, ErrMissingLabelIDField
	}

	if labelIDField.Kind() != reflect.Uint32 {
		return id, labelID, fmt.Errorf("'LabelID' field is %v, expected uint32", labelIDField.Kind())
	}

	labelID = uint32(labelIDField.Uint())

	return id, labelID, nil
}

func setID(e any, id uint32) {
	if _, _, err := checkEntity(e); err != nil {
		panic(err)
	}

	v := reflect.ValueOf(e).Elem()
	idField := v.FieldByName(IDFieldName)
	idField.SetUint(uint64(id))
}
