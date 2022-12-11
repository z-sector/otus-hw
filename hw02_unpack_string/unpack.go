package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const EscapeSymbol = '\\'

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	var (
		builder  strings.Builder
		prev     string
		isEscape bool
	)

	for _, char := range s {
		if isEscape && !(unicode.IsDigit(char) || char == EscapeSymbol) {
			return "", ErrInvalidString
		}

		if char == EscapeSymbol && !isEscape {
			isEscape = true
			continue
		}

		if unicode.IsDigit(char) && !isEscape {
			if prev == "" {
				return "", ErrInvalidString
			}
			count, _ := strconv.Atoi(string(char)) // guaranteed by checking "IsDigit"
			builder.WriteString(strings.Repeat(prev, count))
			prev = ""
			continue
		}

		builder.WriteString(prev)
		prev = string(char)
		isEscape = false
	}

	if isEscape {
		return "", ErrInvalidString
	}
	builder.WriteString(prev)

	return builder.String(), nil
}
