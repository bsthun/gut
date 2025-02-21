package gut

import (
	"errors"
	"reflect"
)

func StructSize(x any) int {
	v := reflect.ValueOf(x)
	return v.NumField()
}

func StructIndex[T any](x any, index int) T {
	v := reflect.ValueOf(x)
	return v.Field(index).Interface().(T)
}

func StructClone(src interface{}, dst interface{}) error {
	srcVal := reflect.ValueOf(src).Elem()
	dstVal := reflect.ValueOf(dst).Elem()

	if srcVal.Kind() != reflect.Struct || dstVal.Kind() != reflect.Struct {
		return errors.New("StructClone: only structs are supported")
	}

	for i := 0; i < srcVal.NumField(); i++ {
		field := srcVal.Type().Field(i)
		srcField := srcVal.Field(i)
		dstField := dstVal.FieldByName(field.Name)

		if dstField.IsValid() && dstField.CanSet() && srcField.Type() == dstField.Type() {
			dstField.Set(srcField)
		}
	}
	return nil
}
