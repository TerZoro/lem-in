package simulation

import (
	"bufio"
	errhandle "lemin/pkg/errors"
	"os"
	"strings"
)

func FileReader(name string) ([]string, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, errhandle.FormatError(errhandle.ErrFileOpen, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, errhandle.FormatError(errhandle.ErrFileRead, err)
	}
	return lines, nil
}

func CheckStartEnd(lines []string) error {
	hasStart := false
	hasEnd := false
	for _, line := range lines {
		if strings.TrimSpace(line) == "##start" {
			hasStart = true
		}
		if strings.TrimSpace(line) == "##end" {
			hasEnd = true
		}
	}
	if !hasStart || !hasEnd {
		return errhandle.FormatError(errhandle.ErrNoStartEnd)
	}
	return nil
}

func ReadArgs() (string, error) {
	args := os.Args[1:]
	if len(args) < 1 {
		return "", errhandle.FormatError(errhandle.ErrNoTestFile)
	}
	return args[0], nil
}
