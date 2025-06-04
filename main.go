package main

import (
	"fmt"
	"lemin/farm"
	"log"
)

func main() {
	lines := FileReader("test")
	CheckStartEnd(lines)
	// Parse ants count
	ants := farm.AntsNum(lines)

	rooms, err := farm.ParseRooms(lines)
	if err != nil {
		log.Printf("error: Could not parse rooms")
	}

	fmt.Printf("Number of ants: %v\n", ants)
	rooms.Printer()
}
