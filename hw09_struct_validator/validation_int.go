package hw09structvalidator

import (
	"strconv"
	"strings"
)

var validIntFunc = map[string]ValidationF[int]{
	intMinRule: validateIntMin,
	intMaxRule: validateIntMax,
	intInRule:  validateIntIn,
}

func validateIntField(value int, fieldName string, cond []Condition) error {
	return validateTypeField[int](value, fieldName, cond, validIntFunc)
}

func validateIntMin(value int, param string) error {
	m, err := strconv.Atoi(param)
	if err != nil {
		return err
	}
	if value < m {
		return TextErr{Err: ErrValidationIntMin}
	}
	return nil
}

func validateIntMax(value int, param string) error {
	m, err := strconv.Atoi(param)
	if err != nil {
		return err
	}
	if value > m {
		return TextErr{Err: ErrValidationIntMax}
	}
	return nil
}

func validateIntIn(value int, param string) error {
	for _, str := range strings.Split(param, setRuleSeparator) {
		is, err := strconv.Atoi(str)
		if err != nil {
			return err
		}

		if value == is {
			return nil
		}
	}

	return TextErr{Err: ErrValidationIntIn}
}
