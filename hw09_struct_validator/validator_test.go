package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/Tyrqvir/otus_hw/hw09_struct_validator/validationerror"
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

	SliceString struct {
		Strings []string `validate:"len:11"`
	}
)

func TestValidateWithoutError(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: Response{
				Code: 200,
				Body: "{}",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1.0.1",
			},
			expectedErr: nil,
		},
		{
			in: Token{
				Header:    []byte{1, 2, 3},
				Payload:   []byte{3, 2, 1},
				Signature: []byte{2, 2, 2},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:    "1b65d722-3200-11ed-a261-0242ac120002",
				Name:  "Vitaly",
				Age:   33,
				Email: "test@gmail.com",
				Role:  "admin",
				Phones: []string{
					"11111111111",
					"22222222222",
				},
				meta: json.RawMessage("{}"),
			},
			expectedErr: nil,
		},
		{
			in: SliceString{
				Strings: []string{
					"11111111111",
					"11111111111",
				},
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			e := Validate(tt.in)

			require.True(t, errors.Is(e, tt.expectedErr))
			_ = tt
		})
	}
}

func TestValidateWithError(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr []string
	}{
		{
			in: Response{
				Code: 201,
				Body: "{}",
			},
			expectedErr: []string{
				"field: [Code] - value not in : [200,404,500], taken - [201]",
			},
		},
		{
			in: User{
				ID:    "1b65d722-3200-11ed-a261-0242ac12001",
				Name:  "Vitaly",
				Age:   5,
				Email: "test@gmail.com",
				Role:  "admin",
				Phones: []string{
					"1111111111",
					"22222222222",
				},
				meta: json.RawMessage("{}"),
			},
			expectedErr: []string{
				"field: [ID] - string must be equal : [36], taken - [1b65d722-3200-11ed-a261-0242ac12001]",
				"field: [Phones[0]] - string must be equal : [11], taken - [1111111111]",
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			e := Validate(tt.in)

			validationErrors := unwrapErrors(e)

			for i, validationError := range validationErrors {
				require.EqualError(t, validationError, tt.expectedErr[i])
			}
			_ = tt
		})
	}
}

func unwrapErrors(e error) validationerror.ValidationErrors {
	var errs validationerror.ValidationErrors
	if !errors.As(e, &errs) {
		return nil
	}
	return errs
}
