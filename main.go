package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("test.txt")
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
	// Parse ants count
	ants, err := strconv.Atoi(lines[0])
	if err != nil || ants <= 0 {
		log.Fatalf("ERROR: invalid number of ants")
	}

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

	var rooms []Room
	for i := 2; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, "#") {
			break
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 3 {
			id, _ := strconv.Atoi(parts[0])
			x, _ := strconv.Atoi(parts[1])
			y, _ := strconv.Atoi(parts[2])
			room := Room{id: id, x: x, y: y}
			rooms = append(rooms, room)
		}
	}

	fmt.Println("Parsed rooms:")
	for _, r := range rooms {
		fmt.Printf("Room: %d (%d, %d)\n", r.id, r.x, r.y)
	}
}
