package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	//nolint:depguard
	"github.com/pkg/errors"
)

const (
	validatePrefixIn    = "in:"
	validatePrefixLen   = "len:"
	validatePrefixRegex = "regexp:"
	validatePrefixMin   = "min:"
	validatePrefixMax   = "max:"
)

var (
	ErrValidationIn    = errors.New("value is not allowed")
	ErrValidationLen   = errors.New("value length does not match expected")
	ErrValidationRegex = errors.New("value does not match regex")
	ErrValidationMin   = errors.New("value is less than min")
	ErrValidationMax   = errors.New("value is larger than max")

	ErrValidationIsNotStruct        = errors.New("value must be a struct")
	ErrValidationForFieldType       = errors.New("unsupported validation for field type")
	ErrValidationUnsupportedTag     = errors.New("unsupported validation tag")
	ErrValidationUnexpectedValueTag = errors.New("value tag has unexpected format")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidatableField interface {
	~int | ~string
}

type ValidationErrors []ValidationError

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Err)
}

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}

	if len(v) == 1 {
		return v[0].Error()
	}

	sb := strings.Builder{}
	for _, err := range v {
		sb.WriteString(fmt.Sprintf("%s\n", err.Error()))
	}

	return sb.String()
}

func Validate(v interface{}) error {
	reflectV := reflect.ValueOf(v)
	if reflectV.Kind() != reflect.Struct {
		return ErrValidationIsNotStruct
	}

	return validateStruct(reflectV)
}

func validateStruct(reflectV reflect.Value) error {
	validationErrors := make(ValidationErrors, 0)
	for i := 0; i < reflectV.Type().NumField(); i++ {
		err := validateField(reflectV.Type().Field(i), reflectV.Field(i))
		if err != nil {
			var validationError ValidationError
			ok := errors.As(err, &validationError)

			if ok {
				validationErrors = append(validationErrors, validationError)
			} else {
				return err
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateField(reflectT reflect.StructField, reflectV reflect.Value) error {
	validateTag := reflectT.Tag.Get("validate")
	if validateTag == "" {
		return nil
	}

	return validateKind(reflectT.Name, reflectV, strings.Split(validateTag, "|"))
}

func validateKind(filedName string, reflectV reflect.Value, validateTags []string) error {
	switch {
	case reflectV.Kind() == reflect.Slice:
		return validateKindSlice(filedName, reflectV.Interface(), validateTags)
	case reflectV.Kind() == reflect.Int:
		return validateKindField(filedName, int(reflectV.Int()), validateTags)
	case reflectV.Kind() == reflect.String:
		return validateKindField(filedName, reflectV.String(), validateTags)
	}

	return nil
}

func validateKindSlice(fieldName string, valueAny any, validateTags []string) error {
	switch value := valueAny.(type) {
	case []int:
		for _, v := range value {
			err := validateKindField(fieldName, v, validateTags)
			if err != nil {
				return err
			}
		}
	case []string:
		for _, v := range value {
			err := validateKindField(fieldName, v, validateTags)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func validateKindField[T ValidatableField](fieldName string, value T, validateTags []string) error {
	valueKind := reflect.ValueOf(value).Kind()

	for _, validateTag := range validateTags {
		if (strings.HasPrefix(validateTag, validatePrefixLen) ||
			strings.HasPrefix(validateTag, validatePrefixRegex)) && valueKind != reflect.String {
			return ValidationError{fieldName, ErrValidationForFieldType}
		}

		if (strings.HasPrefix(validateTag, validatePrefixMin) ||
			strings.HasPrefix(validateTag, validatePrefixMax)) && valueKind != reflect.Int {
			return ValidationError{fieldName, ErrValidationForFieldType}
		}

		var err error
		switch {
		case strings.HasPrefix(validateTag, validatePrefixIn):
			err = validateKindFieldIn(fieldName, value, validateTag)
		case strings.HasPrefix(validateTag, validatePrefixLen):
			err = validateKindFieldLen(fieldName, fmt.Sprint(value), validateTag)
		case strings.HasPrefix(validateTag, validatePrefixRegex):
			err = validateKindFieldRegex(fieldName, fmt.Sprint(value), validateTag)
		case strings.HasPrefix(validateTag, validatePrefixMin):
			intValue, _ := strconv.Atoi(fmt.Sprint(value))
			err = validateKindFieldMin(fieldName, intValue, validateTag)
		case strings.HasPrefix(validateTag, validatePrefixMax):
			intValue, _ := strconv.Atoi(fmt.Sprint(value))
			err = validateKindFieldMax(fieldName, intValue, validateTag)
		default:
			err = ValidationError{fieldName, ErrValidationUnsupportedTag}
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func validateKindFieldIn[T ValidatableField](fieldName string, value T, validateTag string) error {
	allowedValues := strings.Split(strings.TrimPrefix(validateTag, validatePrefixIn), ",")
	if !slices.Contains(allowedValues, fmt.Sprint(value)) {
		return ValidationError{fieldName, ErrValidationIn}
	}

	return nil
}

func validateKindFieldLen(fieldName string, value string, validateTag string) error {
	validationLen, err := strconv.Atoi(strings.TrimPrefix(validateTag, validatePrefixLen))
	if err != nil {
		return ValidationError{fieldName, ErrValidationUnexpectedValueTag}
	}

	if len(value) != validationLen {
		return ValidationError{fieldName, ErrValidationLen}
	}

	return nil
}

func validateKindFieldRegex(fieldName string, value string, validateTag string) error {
	matched, err := regexp.Match(strings.TrimPrefix(validateTag, validatePrefixRegex), []byte(value))
	if err != nil {
		return ValidationError{fieldName, ErrValidationUnexpectedValueTag}
	}

	if !matched {
		return ValidationError{fieldName, ErrValidationRegex}
	}

	return nil
}

func validateKindFieldMin(fieldName string, value int, validateTag string) error {
	validationMin, err := strconv.Atoi(strings.TrimPrefix(validateTag, validatePrefixMin))
	if err != nil {
		return ValidationError{fieldName, ErrValidationUnexpectedValueTag}
	}

	if value < validationMin {
		return ValidationError{fieldName, ErrValidationMin}
	}

	return nil
}

func validateKindFieldMax(fieldName string, value int, validateTag string) error {
	validationMax, err := strconv.Atoi(strings.TrimPrefix(validateTag, validatePrefixMax))
	if err != nil {
		return ValidationError{fieldName, ErrValidationUnexpectedValueTag}
	}

	if value > validationMax {
		return ValidationError{fieldName, ErrValidationMax}
	}

	return nil
}
