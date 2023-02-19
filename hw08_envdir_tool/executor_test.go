package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	type testCase struct {
		cmd        []string
		env        Environment
		returnCode int
	}

	testCases := map[string]testCase{
		"empty cmd": {
			cmd:        []string{},
			env:        Environment{},
			returnCode: 1,
		},
		"cmd without flags": {
			cmd:        []string{"echo"},
			env:        Environment{},
			returnCode: 0,
		},
		"cmd with flag": {
			cmd:        []string{"echo", "-n"},
			env:        Environment{},
			returnCode: 0,
		},
		"direct exit code": {
			cmd:        []string{"sh", "-c", "exit 10"},
			env:        Environment{},
			returnCode: 10,
		},
		"exit code from environment variable": {
			cmd:        []string{"sh", "-c", "exit ${EXIT_CODE}"},
			env:        Environment{"EXIT_CODE": EnvValue{Value: "101"}},
			returnCode: 101,
		},
	}
	for key := range testCases {
		tc := testCases[key]
		name := key
		t.Run(name, func(t *testing.T) {
			actual := RunCmd(tc.cmd, tc.env)
			require.Equal(t, tc.returnCode, actual)
		})
	}

	t.Run("remove env var", func(t *testing.T) {
		varKey, varValue := "TEST_REMOVE_KEY", "TEST_REMOVE_VALUE"
		err := os.Setenv(varKey, varValue)
		require.NoError(t, err)
		defer func() {
			errD := os.Unsetenv(varKey)
			require.NoError(t, errD)
		}()

		env := Environment{
			varKey: EnvValue{NeedRemove: true},
		}
		actual := RunCmd([]string{"echo", "-n"}, env)
		require.Zero(t, actual)
		_, ok := os.LookupEnv(varKey)
		require.False(t, ok)
	})

	t.Run("add env var", func(t *testing.T) {
		varKey, varValue := "TEST_ADD_KEY", "TEST_ADD_VALUE"
		_, ok := os.LookupEnv(varKey)
		require.False(t, ok)

		env := Environment{
			varKey: EnvValue{Value: varValue},
		}
		actual := RunCmd([]string{"echo", "-n"}, env)
		require.Zero(t, actual)
		_, ok = os.LookupEnv(varKey)
		require.True(t, ok)
	})

	t.Run("change env var", func(t *testing.T) {
		varKey, varValue := "TEST_CHANGE_KEY", "TEST_CHANGE_VALUE"
		err := os.Setenv(varKey, varValue)
		require.NoError(t, err)

		newVarValue := "NEW_" + varValue
		env := Environment{
			varKey: EnvValue{Value: newVarValue},
		}
		actual := RunCmd([]string{"echo", "-n"}, env)
		require.Zero(t, actual)
		actualVarValue, ok := os.LookupEnv(varKey)
		require.True(t, ok)
		require.Equal(t, newVarValue, actualVarValue)
	})
}
