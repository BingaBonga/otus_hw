package main

import (
	"bytes"
	"io"
	"os"
	"path"
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
		if !entry.IsDir() && !strings.Contains(entry.Name(), "=") {
			file, err := os.Open(path.Join(dir, entry.Name()))
			if err != nil {
				return nil, err
			}

			defer file.Close()
			fileContentByte, err := io.ReadAll(file)
			if err != nil {
				return nil, err
			}

			fileContent := string(bytes.ReplaceAll(bytes.Split(fileContentByte, []byte("\n"))[0], []byte("\x00"), []byte("\n")))
			env[entry.Name()] = EnvValue{strings.TrimRight(fileContent, " \t\r\n"), len(fileContent) == 0}
		}
	}

	return env, nil
}
