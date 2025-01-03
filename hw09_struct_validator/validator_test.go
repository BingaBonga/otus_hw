package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	//nolint:depguard
	"github.com/google/uuid"
	//nolint:depguard
	"github.com/stretchr/testify/assert"
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

	ValidateTagErrors struct {
		UnsupportedTag           string `validate:"unsupported:200"`
		UnsupportedTagForString1 string `validate:"min:200"`
		UnsupportedTagForString2 string `validate:"max:200"`
		UnsupportedTagForInt1    int    `validate:"len:200"`
		UnsupportedTagForInt2    int    `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		TagWithUnexpectedFormat1 string `validate:"len:dwdswdsw"`
		TagWithUnexpectedFormat2 int    `validate:"min:dwdswdsw"`
		TagWithUnexpectedFormat3 int    `validate:"max:dwdswdsw"`
	}
)

//nolint:lll
func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          App{Version: "12345"},
			expectedErr: nil,
		},
		{
			in:          App{Version: "huawei.com"},
			expectedErr: ValidationError{"Version", ErrValidationLen},
		},
		{
			in:          "StringValue",
			expectedErr: ErrValidationIsNotStruct,
		},
		{
			in:          User{uuid.New().String(), "Something", 22, "email@test.com", UserRole("admin"), []string{"89991111111"}, json.RawMessage{}},
			expectedErr: nil,
		},
		{
			in:          User{uuid.New().String(), "Something", 11, "email@test.com", UserRole("admin"), []string{"89991111111"}, json.RawMessage{}},
			expectedErr: ValidationError{"Age", ErrValidationMin},
		},
		{
			in:          User{uuid.New().String(), "Something", 111, "email@test.com", UserRole("admin"), []string{"89991111111"}, json.RawMessage{}},
			expectedErr: ValidationError{"Age", ErrValidationMax},
		},
		{
			in:          User{uuid.New().String(), "Something", 22, "email", UserRole("admin"), []string{"89991111111"}, json.RawMessage{}},
			expectedErr: ValidationError{"Email", ErrValidationRegex},
		},
		{
			in:          User{uuid.New().String(), "Something", 22, "email@test.com", UserRole("nachalnika"), []string{"89991111111"}, json.RawMessage{}},
			expectedErr: ValidationError{"Role", ErrValidationIn},
		},
		{
			in:          User{uuid.New().String(), "Something", 22, "email@test.com", UserRole("stuff"), []string{"89991111111", "12345"}, json.RawMessage{}},
			expectedErr: ValidationError{"Phones", ErrValidationLen},
		},
		{
			in: User{"", "", 1, "", UserRole(""), []string{"", "89991111111"}, json.RawMessage{}},
			expectedErr: ValidationErrors{
				ValidationError{"ID", ErrValidationLen},
				ValidationError{"Age", ErrValidationMin},
				ValidationError{"Email", ErrValidationRegex},
				ValidationError{"Role", ErrValidationIn},
				ValidationError{"Phones", ErrValidationLen},
			},
		},
		{
			in:          Token{Header: []byte("Something"), Payload: []byte("Something")},
			expectedErr: nil,
		},
		{
			in:          Response{Code: 200, Body: "Something"},
			expectedErr: nil,
		},
		{
			in:          Response{Code: 403, Body: "Something"},
			expectedErr: ValidationError{"Code", ErrValidationIn},
		},
		{
			in: ValidateTagErrors{"", "", "", 0, 0, "", 0, 0},
			expectedErr: ValidationErrors{
				ValidationError{"UnsupportedTag", ErrValidationUnsupportedTag},
				ValidationError{"UnsupportedTagForString1", ErrValidationForFieldType},
				ValidationError{"UnsupportedTagForString2", ErrValidationForFieldType},
				ValidationError{"UnsupportedTagForInt1", ErrValidationForFieldType},
				ValidationError{"UnsupportedTagForInt2", ErrValidationForFieldType},
				ValidationError{"TagWithUnexpectedFormat1", ErrValidationUnexpectedValueTag},
				ValidationError{"TagWithUnexpectedFormat2", ErrValidationUnexpectedValueTag},
				ValidationError{"TagWithUnexpectedFormat3", ErrValidationUnexpectedValueTag},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			var validationErrors ValidationErrors
			ok := errors.As(err, &validationErrors)
			if ok {
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.Equal(t, tt.expectedErr, err)
			}
		})
	}
}
