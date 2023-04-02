package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Nested struct {
		User     User `validate:"nested"`
		Intfield int  `validate:"in:200,404"`
		App      App
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "7f0e3265-ca96-4b33-8858-fef9696cc71b",
				Name:   "Name",
				Age:    30,
				Email:  "somemail@gmail.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
				meta:   []byte{12},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "f",
				Name:   "Name",
				Age:    1,
				Email:  "somemailgmail.com",
				Role:   "role",
				Phones: []string{"0123456789A", "01234567890123456789"},
				meta:   []byte{12},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: ".ID", Err: ErrorArray{TextErr{Err: ErrValidationStrLen}}},
				ValidationError{Field: ".Age", Err: ErrorArray{TextErr{Err: ErrValidationIntMin}}},
				ValidationError{Field: ".Email", Err: ErrorArray{TextErr{Err: ErrValidationStrRegexp}}},
				ValidationError{Field: ".Role", Err: ErrorArray{TextErr{Err: ErrValidationStrIn}}},
				ValidationError{Field: ".Phones[1]", Err: ErrorArray{TextErr{Err: ErrValidationStrLen}}},
			},
		},
		{
			in: App{
				Version: "debug",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "release",
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: ".Version", Err: ErrorArray{TextErr{Err: ErrValidationStrLen}}},
			},
		},
		{
			in: Token{
				Header:    []byte{1, 2, 3},
				Payload:   []byte{1, 2, 3},
				Signature: []byte{1, 2, 3},
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 213,
				Body: "body",
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: ".Code", Err: ErrorArray{TextErr{Err: ErrValidationIntIn}}},
			},
		},
		{
			in: Nested{
				User: User{
					ID:     "7f0e3265-ca96-4b33-8858-fef9696cc71b",
					Name:   "Name",
					Age:    30,
					Email:  "error",
					Role:   "admin",
					Phones: []string{"12345678901"},
					meta:   []byte{12},
				},
				Intfield: 100,
				App: App{
					Version: "release",
				},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: ".User.Email", Err: ErrorArray{TextErr{Err: ErrValidationStrRegexp}}},
				ValidationError{Field: ".Intfield", Err: ErrorArray{TextErr{Err: ErrValidationIntIn}}},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			require.Equal(t, tt.expectedErr, Validate(tt.in))
			_ = tt
		})
	}
}
