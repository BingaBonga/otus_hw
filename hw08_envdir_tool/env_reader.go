package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
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
	dirEntry, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment, len(dirEntry))
	for _, entry := range dirEntry {
		if !entry.IsDir() {
			file, err := os.Open(fmt.Sprintf("%s/%s", dir, entry.Name()))
			if err != nil {
				return nil, err
			}

			defer file.Close()
			fileContentByte, err := io.ReadAll(file)
			if err != nil {
				return nil, err
			}

			fileContent := strings.TrimRight(string(bytes.Replace(fileContentByte, []byte("\x00"), []byte("\n"), -1)), " \t\r\n")
			env[entry.Name()] = EnvValue{fileContent, len(fileContent) == 0}
		}
	}

	return env, nil
}
