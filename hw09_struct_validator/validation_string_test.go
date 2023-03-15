package hw09structvalidator

import (
	"errors"
	"fmt"
	"regexp/syntax"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateStringIn(t *testing.T) {
	tests := []struct {
		value       string
		param       string
		expectedErr error
	}{
		{
			value:       "aaa",
			param:       "aaa,bbb",
			expectedErr: nil,
		},
		{
			value:       "bbb",
			param:       "aaa,bbb",
			expectedErr: nil,
		},
		{
			value:       "aaa",
			param:       "",
			expectedErr: TextErr{Err: ErrValidationStrIn},
		},
		{
			value:       "ccc",
			param:       "aaa,bbb",
			expectedErr: TextErr{Err: ErrValidationStrIn},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tc := tc
			t.Parallel()

			actErr := validateStringIn(tc.value, tc.param)

			require.IsType(t, tc.expectedErr, actErr)
			if errors.As(actErr, &TextErr{}) {
				require.Equal(t, tc.expectedErr, actErr)
			}

			_ = tc
		})
	}
}

func TestValidateStringRegex(t *testing.T) {
	tests := []struct {
		value       string
		param       string
		expectedErr error
	}{
		{
			value:       "aaa",
			param:       "[",
			expectedErr: &syntax.Error{},
		},
		{
			value:       "",
			param:       "^\\d+$",
			expectedErr: TextErr{Err: ErrValidationStrRegexp},
		},
		{
			value:       "foo",
			param:       "",
			expectedErr: nil,
		},
		{
			value:       "123456",
			param:       "^\\d+$",
			expectedErr: nil,
		},
		{
			value:       "error case",
			param:       "foo",
			expectedErr: TextErr{Err: ErrValidationStrRegexp},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tc := tc
			t.Parallel()

			actErr := validateStringRegexp(tc.value, tc.param)

			require.IsType(t, tc.expectedErr, actErr)
			if errors.As(actErr, &TextErr{}) {
				require.Equal(t, tc.expectedErr, actErr)
			}

			_ = tc
		})
	}
}

func TestValidateStringLen(t *testing.T) {
	tests := []struct {
		value       string
		param       string
		expectedErr error
	}{
		{
			value:       "",
			param:       "",
			expectedErr: &strconv.NumError{},
		},
		{
			value:       "",
			param:       "asd",
			expectedErr: &strconv.NumError{},
		},
		{
			value:       "",
			param:       "0",
			expectedErr: nil,
		},
		{
			value:       "string",
			param:       "6",
			expectedErr: nil,
		},
		{
			value:       "string",
			param:       "5",
			expectedErr: TextErr{Err: ErrValidationStrLen},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tc := tc
			t.Parallel()

			actErr := validateStringLen(tc.value, tc.param)

			require.IsType(t, tc.expectedErr, actErr)
			if errors.As(actErr, &TextErr{}) {
				require.Equal(t, tc.expectedErr, actErr)
			}

			_ = tc
		})
	}
}
