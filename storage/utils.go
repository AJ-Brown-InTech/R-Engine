package storage

import (
	"fmt"
	"reflect"

)

// ObjectToDatabase converts a struct to column names and values.
func ObjectToDatabase(obj interface{}) ([]string, []interface{}, error) {
	var columns []string
	var values []interface{}

	objValue := reflect.ValueOf(obj)
	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}

	if objValue.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("input must be a struct or a pointer to a struct")
	}

	objType := objValue.Type()
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		tag := field.Tag.Get("db")
		if tag != "" && tag != "-" {
			fieldValue := objValue.Field(i)
			if fieldValue.Kind() != reflect.Ptr || !fieldValue.IsNil() {
				columns = append(columns, tag)
				values = append(values, fieldValue.Interface())
			}
		}
	}
	return columns, values, nil
}
