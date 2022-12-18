package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
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

	isDigit := func(c rune) bool { return c >= '0' && c <= '9' }

	for _, char := range s {
		if isEscape && !(isDigit(char) || char == EscapeSymbol) {
			return "", ErrInvalidString
		}

		if char == EscapeSymbol && !isEscape {
			isEscape = true
			continue
		}

		if isDigit(char) && !isEscape {
			if prev == "" {
				return "", ErrInvalidString
			}
			count, _ := strconv.Atoi(string(char)) // guaranteed by checking "isDigit"
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
