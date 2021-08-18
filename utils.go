package main

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
)

// exists returns whether the given file or directory exists
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// readFile reads the given file and returns its string content.
func readFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	text := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text += scanner.Text() + "\n"
	}

	return text, nil
}

// getApolloDir returns the absolute path of the apollo directory, makes sure it exists.
func getApolloDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	dir = filepath.Join(dir, ".apollo")

	if !exists(dir) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			panic(err)
		}
	}

	return dir
}

// readJSON reads the given json string and returns a Configuration struct.
func readJSON(jsonStr string) (*Configuration, error) {
	config := Configuration{}
	err := json.Unmarshal([]byte(jsonStr), &config)
	return &config, err
}

// writeFile writes the given string to the given file.
func writeFile(path, text string) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(text)
}

// jsonify accepts a Configuration object and converts it to JSON.
func jsonify(config *Configuration) (string, error) {
	jsonBytes, err := json.MarshalIndent(config, "", "    ")
	return string(jsonBytes), err
}
