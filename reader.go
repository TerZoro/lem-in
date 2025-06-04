package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func FileReader(name string) []string {
	file, err := os.Open(name + ".txt")
	if err != nil {
		log.Fatalf("ERROR: could not open file: %v", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("ERROR: reading file: %v", err)
	}
	return lines
}

func CheckStartEnd(lines []string) {
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
		log.Fatalf("ERROR: missing ##start or ##end")
	}
}
