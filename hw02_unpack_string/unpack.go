package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var buffer *rune
	asterisk := false
	output := strings.Builder{}

	for _, v := range input {
		// copy loop var
		char := v
		if digit, err := strconv.Atoi(string(char)); err == nil && !asterisk {
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

		switch {
		case asterisk:
			buffer = &char
			asterisk = false
		case char == '\\':
			asterisk = true
		default:
			buffer = &char
		}
	}

	if buffer != nil {
		output.WriteRune(*buffer)
	}
	return output.String(), nil
}
