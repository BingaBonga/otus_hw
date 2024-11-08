package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var buffer *rune
	slash := false
	output := strings.Builder{}

	for _, v := range input {
		// copy loop var
		char := v

		// char to digit convert
		digit, errDigitConvert := strconv.Atoi(string(char))
		if errDigitConvert == nil && !slash {
			if buffer == nil {
				return "", ErrInvalidString
			}

			output.WriteString(strings.Repeat(string(*buffer), digit))
			buffer = nil
			continue
		} else if buffer != nil {
			output.WriteRune(*buffer)
			buffer = nil
		}

		// slash can be before slash or digit
		if slash && char != '\\' && errDigitConvert != nil {
			return "", ErrInvalidString
		}

		switch {
		case slash:
			buffer = &char
			slash = false
		case char == '\\':
			slash = true
		default:
			buffer = &char
		}
	}

	// check on last slash
	if slash {
		return "", ErrInvalidString
	}

	// don't forget last char
	if buffer != nil {
		output.WriteRune(*buffer)
	}
	return output.String(), nil
}
