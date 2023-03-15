package hw09structvalidator

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateIntMin(t *testing.T) {
	tests := []struct {
		value       int
		param       string
		isErr       bool
		expectedErr error
	}{
		{
			value:       0,
			param:       "foo",
			isErr:       true,
			expectedErr: &strconv.NumError{},
		},
		{
			value:       2,
			param:       "2",
			expectedErr: nil,
		},
		{
			value:       3,
			param:       "2",
			expectedErr: nil,
		},
		{
			value:       2,
			param:       "3",
			expectedErr: TextErr{Err: ErrValidationIntMin},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tc := tc
			t.Parallel()

			actErr := validateIntMin(tc.value, tc.param)

			require.IsType(t, tc.expectedErr, actErr)
			if errors.As(actErr, &TextErr{}) {
				require.Equal(t, tc.expectedErr, actErr)
			}

			_ = tc
		})
	}
}

func TestValidateIntMax(t *testing.T) {
	tests := []struct {
		value       int
		param       string
		expectedErr error
	}{
		{
			value:       0,
			param:       "foo",
			expectedErr: &strconv.NumError{},
		},
		{
			value:       2,
			param:       "2",
			expectedErr: nil,
		},
		{
			value:       1,
			param:       "2",
			expectedErr: nil,
		},
		{
			value:       4,
			param:       "3",
			expectedErr: TextErr{Err: ErrValidationIntMax},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tc := tc
			t.Parallel()

			actErr := validateIntMax(tc.value, tc.param)

			require.IsType(t, tc.expectedErr, actErr)
			if errors.As(actErr, &TextErr{}) {
				require.Equal(t, tc.expectedErr, actErr)
			}

			_ = tc
		})
	}
}

func TestValidateIntIn(t *testing.T) {
	tests := []struct {
		value       int
		param       string
		expectedErr error
	}{
		{
			value:       0,
			param:       "foo",
			expectedErr: &strconv.NumError{},
		},
		{
			value:       2,
			param:       "2,4",
			expectedErr: nil,
		},
		{
			value:       4,
			param:       "2,4",
			expectedErr: nil,
		},
		{
			value:       1,
			param:       "2,4",
			expectedErr: TextErr{Err: ErrValidationIntIn},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tc := tc
			t.Parallel()

			actErr := validateIntIn(tc.value, tc.param)

			require.IsType(t, tc.expectedErr, actErr)
			if errors.As(actErr, &TextErr{}) {
				require.Equal(t, tc.expectedErr, actErr)
			}

			_ = tc
		})
	}
}
