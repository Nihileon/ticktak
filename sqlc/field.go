package sqlc

import (
	"errors"
	"reflect"
)

type Fields struct {
	fieldMap   map[string]valType
	fieldSlice []string
}

func NewFieldMap(input interface{}) (fields Fields, err error) {
	value := reflect.ValueOf(input)
	if value.Kind() != reflect.Struct {
		return fields, errors.New("input must be a struct")
	}

	fields.fieldMap = make(map[string]valType)

	base := value.Type()

	for i := 0; i < base.NumField(); i++ {
		field := base.Field(i).Tag.Get("json")
		val := reflect.Indirect(value).Field(i).Interface()
		fields.fieldMap[field] = val
	}
	return fields, nil
}

//
func NewFieldSlice(input interface{}) (fields Fields, err error) {
	value := reflect.ValueOf(input)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return fields, errors.New("input must be a struct")
	}

	base := value.Type()

	for i := 0; i < base.NumField(); i++ {
		fields.fieldSlice = append(fields.fieldSlice, base.Field(i).Tag.Get("json"))
	}
	return fields, nil
}
