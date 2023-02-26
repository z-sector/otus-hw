package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"path"
	"unicode"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	result := make(Environment, len(entries))
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		if info.Size() == 0 {
			result[entry.Name()] = EnvValue{NeedRemove: true}
			continue
		}

		value, err := getValueFromFile(path.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		if value == "" {
			result[entry.Name()] = EnvValue{NeedRemove: true}
			continue
		}

		result[entry.Name()] = EnvValue{Value: value}
	}
	return result, nil
}

func getValueFromFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	reader := bufio.NewReader(file)
	line, err := reader.ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	line = bytes.ReplaceAll(line, []byte{0x00}, []byte("\n"))
	return string(bytes.TrimRightFunc(line, unicode.IsSpace)), nil
}
