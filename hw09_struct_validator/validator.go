package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

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

	ErrValidationTypeWrongTag   = errors.New("field type unsupported validation")
	ErrValidationUnsupportedTag = errors.New("unsupported validation tag")
	ErrValidationValueTag       = errors.New("value tag has unexpected format")
	ErrValidationValueIsNil     = errors.New("value is nil")
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
	if v == nil || len(v) == 0 {
		return ""
	}

	sb := strings.Builder{}
	for _, err := range v {
		sb.WriteString(fmt.Sprintf("%s\n", err.Error()))
	}

	return sb.String()
}

func Validate(v interface{}) error {
	reflectV := reflect.ValueOf(&v)
	if reflectV.Kind() != reflect.Struct {
		return errors.New("value must be a struct or a pointer to struct")
	}

	return validateStruct(reflectV)
}

func validateStruct(reflectV reflect.Value) error {
	validationErrors := make(ValidationErrors, 0)
	for i := 0; i < reflectV.NumField(); i++ {
		reflectT := reflectV.Type()
		err := validateField(reflectT.Field(i), reflectV.Field(i))
		if err != nil {
			var validationError ValidationError
			ok := errors.Is(err, &validationError)

			if ok {
				validationErrors = append(validationErrors, validationError)
			} else {
				return err
			}
		}
	}

	return validationErrors
}

func validateField(reflectT reflect.StructField, reflectV reflect.Value) error {
	validateTags := strings.Split(reflectT.Tag.Get("validate"), "|")
	if len(validateTags) == 0 {
		return nil
	}

	return validateKind(reflectT.Name, reflectV, validateTags)
}

func validateKind(filedName string, reflectV reflect.Value, validateTags []string) error {
	switch {
	case reflectV.Kind() == reflect.Slice:
		return validateKindSlice(filedName, reflectV.Interface(), validateTags)
	case reflectV.Kind() == reflect.String:
	case reflectV.Kind() == reflect.Int:
		switch value := reflectV.Interface().(type) {
		case int:
		case string:
			return validateKindField(filedName, value, validateTags)
		}
	}

	return nil
}

func validateKindSlice(fieldName string, valueAny any, validateTags []string) error {
	switch value := valueAny.(type) {
	case []int:
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
	reflectV := reflect.ValueOf(&value)
	valueKind := reflectV.Kind()

	if reflectV.IsNil() {
		return ValidationError{fieldName, ErrValidationValueIsNil}
	}

	for _, validateTag := range validateTags {
		if (strings.HasPrefix(validateTag, validatePrefixLen) || strings.HasPrefix(validateTag, validatePrefixRegex)) && valueKind != reflect.String {
			return ValidationError{fieldName, errors.Wrap(ErrValidationTypeWrongTag, fmt.Sprintf("for type %v validation tag %s is not allowed", valueKind, validateTag))}
		}

		if (strings.HasPrefix(validateTag, validatePrefixMin) || strings.HasPrefix(validateTag, validatePrefixMax)) && valueKind != reflect.Int {
			return ValidationError{fieldName, errors.Wrap(ErrValidationTypeWrongTag, fmt.Sprintf("for type %v validation tag %s is not allowed", valueKind, validateTag))}
		}

		var err error
		if strings.HasPrefix(validateTag, validatePrefixIn) {
			err = validateKindFieldIn(fieldName, value, validateTag)
		} else if strings.HasPrefix(validateTag, validatePrefixLen) {
			err = validateKindFieldLen(fieldName, fmt.Sprint(value), validateTag)
		} else if strings.HasPrefix(validateTag, validatePrefixRegex) {
			err = validateKindFieldRegex(fieldName, fmt.Sprint(value), validateTag)
		} else if strings.HasPrefix(validateTag, validatePrefixMin) {
			intValue, _ := strconv.Atoi(fmt.Sprint(value))
			err = validateKindFieldMin(fieldName, intValue, validateTag)
		} else if strings.HasPrefix(validateTag, validatePrefixMax) {
			intValue, _ := strconv.Atoi(fmt.Sprint(value))
			err = validateKindFieldMax(fieldName, intValue, validateTag)
		} else {
			err = ValidationError{fieldName, ErrValidationUnsupportedTag}
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func validateKindFieldIn[T ValidatableField](fieldName string, value T, validateTag string) error {
	allowedValues := strings.Split(strings.TrimLeft(validateTag, validatePrefixIn), ",")
	if !slices.Contains(allowedValues, fmt.Sprint(value)) {
		return ValidationError{fieldName, errors.Wrap(ErrValidationIn, fmt.Sprintf("allowed values: %v, current value is %s", allowedValues, fmt.Sprint(value)))}
	}

	return nil
}

func validateKindFieldLen(fieldName string, value string, validateTag string) error {
	validationLen, err := strconv.Atoi(strings.TrimLeft(validateTag, validatePrefixLen))
	if err != nil {
		return ValidationError{fieldName, errors.Wrap(ErrValidationValueTag, fmt.Sprintf("validation value must be a number, current value %s", strings.TrimLeft(validateTag, validatePrefixLen)))}
	}

	if len(value) != validationLen {
		return ValidationError{fieldName, errors.Wrap(ErrValidationLen, fmt.Sprintf("value length must be %v, current lenght %v", validationLen, len(value)))}
	}

	return nil
}

func validateKindFieldRegex(fieldName string, value string, validateTag string) error {
	matched, err := regexp.Match(strings.TrimLeft(validateTag, validatePrefixRegex), []byte(value))
	if err != nil {
		return ValidationError{fieldName, errors.Wrap(ErrValidationValueTag, fmt.Sprintf("validation value must be a regex, current value %s", strings.TrimLeft(validateTag, validatePrefixRegex)))}
	}

	if !matched {
		return ValidationError{fieldName, errors.Wrap(ErrValidationRegex, fmt.Sprintf("value must be macthed %s, current value %s", strings.TrimLeft(validateTag, validatePrefixRegex), value))}
	}

	return nil
}

func validateKindFieldMin(fieldName string, value int, validateTag string) error {
	validationMin, err := strconv.Atoi(strings.TrimLeft(validateTag, validatePrefixMin))
	if err != nil {
		return ValidationError{fieldName, errors.Wrap(ErrValidationValueTag, fmt.Sprintf("validation value must be a number, current value %s", strings.TrimLeft(validateTag, validatePrefixMin)))}
	}

	if value < validationMin {
		return ValidationError{fieldName, errors.Wrap(ErrValidationMin, fmt.Sprintf("value must be less %v, current value %v", validationMin, value))}
	}

	return nil
}

func validateKindFieldMax(fieldName string, value int, validateTag string) error {
	validationMax, err := strconv.Atoi(strings.TrimLeft(validateTag, validatePrefixMax))
	if err != nil {
		return ValidationError{fieldName, errors.Wrap(ErrValidationValueTag, fmt.Sprintf("validation value must be a number, current value %s", strings.TrimLeft(validateTag, validatePrefixMax)))}
	}

	if value > validationMax {
		return ValidationError{fieldName, errors.Wrap(ErrValidationMax, fmt.Sprintf("value musn't larger %v, current value %v", validationMax, value))}
	}

	return nil
}
