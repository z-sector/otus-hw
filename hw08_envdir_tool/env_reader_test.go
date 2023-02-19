package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("read env variables", func(t *testing.T) {
		expected := Environment{
			"BAR":   {Value: "bar", NeedRemove: false},
			"EMPTY": {Value: "", NeedRemove: true},
			"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": {Value: "\"hello\"", NeedRemove: false},
			"UNSET": {Value: "", NeedRemove: true},
		}
		env, err := ReadDir("./testdata/env")
		require.NoError(t, err)
		require.Equal(t, expected, env)
	})

	t.Run("path not found", func(t *testing.T) {
		var expected Environment
		env, err := ReadDir("./testdata/missing_dir")
		require.Error(t, err)
		require.Equal(t, expected, env)
	})

	t.Run("empty dir", func(t *testing.T) {
		expected := Environment{}
		env, err := ReadDir("./testdata/empty_dir")
		require.NoError(t, err)
		require.Equal(t, expected, env)
	})
}
