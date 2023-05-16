package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/eduardluft/repono/pkg/color"
)

func main() {
	flag.Parse()
	file := flag.Arg(0)
	output := entrypoint(file)
	fmt.Print(string(output))
}

func entrypoint(file string) []byte {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		stderr("No file given or found with first argument: \"%s\"", file)
		return nil
	}
	basePath := filepath.Dir(file)

	output := input
	for hasPlaceholder(output) {
		output = process(output, basePath)
	}
	return output
}

func process(data []byte, basePath string) []byte {
	placeholders, err := findPlaceholder(data)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	output := string(data)
	if err != nil {
		fmt.Print("error in reading file.File")
		return nil
	}

	for _, placeholder := range placeholders {
		filename, fileContent := findFileByPlaceholder(placeholder, basePath)

		if isSubPath(filename) {
			fileContentPlaceholders, err := findPlaceholder(fileContent)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			for _, fileContentPlaceholder := range fileContentPlaceholders {
				path := filepath.Dir(filename)
				subPathPlaceholder := extendPlaceholderWithPath(fileContentPlaceholder, path)
				fileContent = []byte(strings.ReplaceAll(string(fileContent), fileContentPlaceholder, subPathPlaceholder))
			}
		}

		output = strings.ReplaceAll(string(output), placeholder, string(fileContent))

	}

	return []byte(output)
}

func extendPlaceholderWithPath(placeholder string, path string) string {
	filename := strings.Trim(placeholder, "###")
	fullFilename := filepath.Join(path, filename)
	return "###" + fullFilename + "###"
}

func findFileByPlaceholder(placeholder string, path string) (string, []byte) {
	filename := strings.ReplaceAll(placeholder, "###", "")
	fullFilename := filepath.Join(path, filename)
	data, err := ioutil.ReadFile(fullFilename)
	if err != nil {
		stderr("Cant find file for placeholder: %s", fullFilename)
		return "", nil
	}

	return filename, data
}

func findPlaceholder(data []byte) ([]string, error) {
	pattern := regexp.MustCompile(`###.*###`)
	matches := pattern.FindAllStringSubmatch(string(data), -1)
	placeholders := []string{}
	for _, match := range matches {
		placeholders = append(placeholders, match[0])
	}

	return placeholders, nil
}

func hasPlaceholder(data []byte) bool {
	pattern := regexp.MustCompile(`###.*###`)
	matches := pattern.FindAllStringSubmatch(string(data), -1)

	return len(matches) != 0
}

func isSubPath(data string) bool {
	pattern := regexp.MustCompile(`.+\/.+`)
	matches := pattern.FindAllStringSubmatch(string(data), -1)

	return len(matches) != 0
}

func stderr(msg string, a ...any) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(color.Red+"Err: "+msg+color.Reset, a...))
}
