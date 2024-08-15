package util

import (
	"fmt"
	"reflect"
)

// ListFields lists the fields of a struct
func ListFields(t any) (r []string, err error) {
	typeOfT := reflect.TypeOf(t)
	if typeOfT.Kind() != reflect.Struct {
		err = fmt.Errorf("can't reflect the fields of non-struct type %s", typeOfT.Elem().Name())
		return
	}

	fields := reflect.VisibleFields(reflect.TypeOf(t))
	for _, f := range fields {
		r = append(r, f.Name)
	}
	return
}

// ListNilFields lists the nil fields of a struct
func ListNilFields(t any) (nilFieldNames []string, err error) {
	typeOfT := reflect.TypeOf(t)
	if typeOfT.Kind() != reflect.Struct {
		err = fmt.Errorf("can't reflect the fields of non-struct type %s", typeOfT.Elem().Name())
		return
	}

	r := reflect.ValueOf(t)
	fields := reflect.VisibleFields(reflect.TypeOf(t))
	for _, f := range fields {
		field := reflect.Indirect(r).FieldByIndex(f.Index)
		if field.IsNil() {
			nilFieldNames = append(nilFieldNames, f.Name)
		}
	}
	return
}
