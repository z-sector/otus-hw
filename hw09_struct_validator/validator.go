package hw09structvalidator

import (
	"errors"
	"reflect"
)

func Validate(v interface{}) error {
	return validateStruct(v, "")
}

func validateStruct(v interface{}, parentName string) error {
	fields, err := parseTags(v, parentName)
	if err != nil {
		return err
	}

	validationErrors := make(ValidationErrors, 0)
	var verrs ValidationErrors
	for _, f := range fields {
		err = validateField(f)
		if err != nil {
			if !errors.As(err, &verrs) {
				return err
			}
			validationErrors = append(validationErrors, verrs...)
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateField(field field) error {
	cond, err := parseConditions(field.tag)
	if err != nil {
		return err
	}

	switch field.kind { //nolint:exhaustive
	case reflect.Struct:
		item := field.value.Interface()
		err = validateStruct(item, field.name)
	case reflect.String:
		return validateStringField(field.value.String(), field.name, cond)
	case reflect.Int:
		err = validateIntField(int(field.value.Int()), field.name, cond)
	case reflect.Slice:
		switch field.value.Type().Elem().Kind() { //nolint:exhaustive
		case reflect.String:
			slice, _ := field.value.Interface().([]string)
			err = validateSliceField[string](slice, field.name, cond, validateStringField)
		case reflect.Int:
			slice, _ := field.value.Interface().([]int)
			err = validateSliceField[int](slice, field.name, cond, validateIntField)
		default:
			err = ErrUnsupported
		}
	default:
		err = ErrUnsupported
	}
	return err
}
