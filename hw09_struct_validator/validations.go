package hw09structvalidator

import (
	"errors"
	"fmt"
)

type (
	ValidationF[T any]     func(v T, param string) error
	ValidTypeFieldF[T any] func(value T, fieldName string, cond []Condition) error
)

func validateSliceField[T any](value []T, fieldName string, cond []Condition, f ValidTypeFieldF[T]) error {
	validationErrors := make(ValidationErrors, 0)
	var verrs ValidationErrors

	for i, v := range value {
		sliceFieldName := fmt.Sprintf("%s[%d]", fieldName, i)
		err := f(v, sliceFieldName, cond)
		if err != nil {
			if !errors.As(err, &verrs) {
				return err
			}
			validationErrors = append(validationErrors, verrs...)
		}
	}

	if len(validationErrors) != 0 {
		return validationErrors
	}

	return nil
}

func validateTypeField[T any](value T, fieldName string, cond []Condition, vFunc map[string]ValidationF[T]) error {
	var err error
	var errText TextErr

	errArr := make(ErrorArray, 0)
	for _, c := range cond {
		f, ok := vFunc[c.operator]
		if !ok {
			return ErrUnknownTag
		}
		err = f(value, c.operand)

		if err != nil {
			if !errors.As(err, &errText) {
				return err
			}
			errArr = append(errArr, errText)
		}
	}

	if len(errArr) > 0 {
		return ValidationErrors{ValidationError{Field: fieldName, Err: errArr}}
	}

	return nil
}
