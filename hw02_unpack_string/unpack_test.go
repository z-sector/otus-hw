package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "abcd", expected: "abcd"},
		{input: "e 2e", expected: "e  e"},
		{input: "x1y1", expected: "xy"},
		{input: "c0", expected: ""},
		{input: "абв", expected: "абв"},
		{input: "а1б2в3", expected: "аббввв"},
		{input: `а\1б2`, expected: "а1бб"},
		{input: `!@#$%^&*()`, expected: "!@#$%^&*()"},
		{input: `!@2#3$`, expected: "!@@###$"},
		{input: `!@\2#\3$`, expected: "!@2#3$"},
		{input: `世2a2界3b3`, expected: "世世aa界界界bbb"},
		{input: `🌀0a1🍣2ф3`, expected: "a🍣🍣ффф"},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: `\11\22\33`, expected: `122333`},
		{input: `\\\\`, expected: `\\`},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", `a\`, "aaa10b", `d2\y`, `\\\`}
	for i := range invalidStrings {
		tc := invalidStrings[i]
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
