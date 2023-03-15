package hw09structvalidator

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrType             = errors.New("invalid type, expected struct")
	ErrValidationFormat = errors.New("invalid validation format, expected [operator[:operand]]")
	ErrUnknownTag       = errors.New("unknown tag")
	ErrUnsupported      = errors.New("unsupported type")

	ErrValidationIntMin = errors.New("value less than expected")
	ErrValidationIntMax = errors.New("value greater than expected")
	ErrValidationIntIn  = errors.New("value doesn't match a subset of int")

	ErrValidationStrLen    = errors.New("string's length is not as expected")
	ErrValidationStrRegexp = errors.New("string doesn't match the regexp")
	ErrValidationStrIn     = errors.New("string doesn't match the substring")
)

type TextErr struct {
	Err error
}

func (t TextErr) Error() string {
	return t.Err.Error()
}

type ErrorArray []TextErr

func (e ErrorArray) Error() string {
	sb := strings.Builder{}
	for _, err := range e {
		sb.WriteString(fmt.Sprintf("%s, ", err.Error()))
	}
	errs := sb.String()
	return strings.TrimSuffix(errs, ", ")
}

type ValidationError struct {
	Field string
	Err   ErrorArray
}

func (v ValidationError) Error() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s: %s", v.Field, v.Err.Error()))
	return sb.String()
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("validation errors, total %d:\n", len(v)))

	for _, err := range v {
		sb.WriteString(fmt.Sprintf("%s\n", err.Error()))
	}

	return sb.String()
}
