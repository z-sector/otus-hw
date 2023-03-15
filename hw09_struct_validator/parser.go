package hw09structvalidator

import (
	"fmt"
	"reflect"
)

type tagT string

type field struct {
	name  string
	value reflect.Value
	tag   tagT
	kind  reflect.Kind
}

func parseTags(v interface{}, parentName string) ([]field, error) {
	valueOf := reflect.ValueOf(v)
	if valueOf.Kind() != reflect.Struct {
		return nil, ErrType
	}

	fields := make([]field, 0)
	for i := 0; i < valueOf.NumField(); i++ {
		fieldValue := valueOf.Field(i)
		structFieldType := valueOf.Type().Field(i)

		if !structFieldType.IsExported() {
			continue
		}

		tagRaw, ok := structFieldType.Tag.Lookup(defaultTagName)
		if !ok {
			continue
		}

		name := structFieldType.Name
		if name != "" {
			name = fmt.Sprintf("%s.%s", parentName, name)
		}

		fields = append(
			fields,
			field{
				name:  name,
				value: fieldValue,
				tag:   tagT(tagRaw),
				kind:  structFieldType.Type.Kind(),
			},
		)
	}

	return fields, nil
}
