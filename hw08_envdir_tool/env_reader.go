package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
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
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)

	for _, file := range dirEntries {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			return env, err
		}

		filePath := filepath.Join(dir, info.Name())
		envName, envValue, err := getEnvByFilePath(filePath, info)

		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}

		env[envName] = envValue
	}

	return env, nil
}

func getEnvByFilePath(filePath string, info fs.FileInfo) (string, EnvValue, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", EnvValue{}, err
	}
	defer file.Close()

	envName := prepareEnvName(info.Name())
	envValue := prepareEnvValue(file, info)

	return envName, envValue, nil
}

func prepareEnvName(name string) string {
	return strings.ReplaceAll(name, "=", "")
}

func prepareEnvValue(file *os.File, info os.FileInfo) EnvValue {
	fileReader := bufio.NewReader(file)
	value, _, _ := fileReader.ReadLine()

	value = bytes.ReplaceAll(value, []byte("\x00"), []byte("\n"))
	value = bytes.TrimRight(value, " ")

	return EnvValue{
		Value:      string(value),
		NeedRemove: info.Size() == 0,
	}
}
