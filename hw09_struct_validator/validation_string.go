package hw09structvalidator

import (
	"regexp"
	"strconv"
	"strings"
)

var validStringFunc = map[string]ValidationF[string]{
	stringRegexpRule: validateStringRegexp,
	stringInRule:     validateStringIn,
	stringLenRule:    validateStringLen,
}

func validateStringField(value string, fieldName string, cond []Condition) error {
	return validateTypeField[string](value, fieldName, cond, validStringFunc)
}

func validateStringRegexp(value string, param string) error {
	re, err := regexp.Compile(param)
	if err != nil {
		return err
	}
	if !re.MatchString(value) {
		return TextErr{Err: ErrValidationStrRegexp}
	}
	return nil
}

func validateStringIn(value string, param string) error {
	if value == "" && param == "" {
		return nil
	}
	for _, is := range strings.Split(param, setRuleSeparator) {
		if value == is {
			return nil
		}
	}
	return TextErr{Err: ErrValidationStrIn}
}

func validateStringLen(value string, param string) error {
	expLen, err := strconv.Atoi(param)
	if err != nil {
		return err
	}
	if len(value) != expLen {
		return TextErr{Err: ErrValidationStrLen}
	}
	return nil
}
