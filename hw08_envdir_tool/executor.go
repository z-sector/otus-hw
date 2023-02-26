package main

import (
	"errors"
	"os"
	"os/exec"
)

const (
	returnCodeOk  = 0
	returnCodeErr = 1
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return returnCodeErr
	}

	if err := cleanEnvs(env); err != nil {
		return returnCodeErr
	}

	name := cmd[0]
	args := cmd[1:]
	c := exec.Command(name, args...)
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout

	err := c.Run()
	if err != nil {
		var exitErr *exec.ExitError

		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return returnCodeErr
	}

	return returnCodeOk
}

func cleanEnvs(env Environment) error {
	for name, val := range env {
		var err error

		if val.NeedRemove {
			err = os.Unsetenv(name)
		} else {
			err = os.Setenv(name, val.Value)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
