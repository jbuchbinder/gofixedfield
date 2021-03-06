package gofixedfield

import (
	"fmt"
	"reflect"
)

func debugStruct(i interface{}) {
	val := reflect.ValueOf(i).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		fmt.Printf("Field Name: %s, Field Value: %v, 'fixed' Tag Value: %s, 'csv' Tag Value: %s, Type: %s\n", typeField.Name, valueField.Interface(), tag.Get("fixed"), tag.Get("csv"), typeField.Type.Kind())
	}
}
