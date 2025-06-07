package main

import (
	"fmt"
	"lemin/farm"
	"log"
)

func main() {
	lines := FileReader("test.txt")
	CheckStartEnd(lines)
	// Parse ants count
	ants := farm.AntsNum(lines)

	rooms, err := farm.ParseRooms(lines)
	if err != nil {
		log.Printf("error: Could not parse rooms")
	}

	// Debug information (comment out for final output)
	/*
		fmt.Printf("Number of ants: %v\n", ants)

		fmt.Println("\n=== Debug Information ===")
		start, _ := rooms.SearchStart()
		end, _ := rooms.SearchEnd()
		if start != nil {
			fmt.Printf("Start room found - ID: %d, Coords: (%d, %d)\n", start.ID, start.X, start.Y)
			fmt.Printf("Start room links to: ")
			for _, link := range start.Links {
				fmt.Printf("%d ", link.ID)
			}
			fmt.Println()
		} else {
			fmt.Println("WARNING: No start room found!")
		}

		if end != nil {
			fmt.Printf("End room found - ID: %d, Coords: (%d, %d)\n", end.ID, end.X, end.Y)
			fmt.Printf("End room is linked from: ")
			for _, room := range rooms.List {
				for _, link := range room.Links {
					if link.ID == end.ID {
						fmt.Printf("%d ", room.ID)
					}
				}
			}
			fmt.Println()
		} else {
			fmt.Println("WARNING: No end room found!")
		}

		fmt.Println("\nAll rooms and their links:")
		for _, room := range rooms.List {
			fmt.Printf("Room %d links to: ", room.ID)
			for _, link := range room.Links {
				fmt.Printf("%d ", link.ID)
			}
			fmt.Println()
		}
		fmt.Println("=== End Debug Info ===\n")

		rooms.Printer()

		fmt.Printf("\n=== Found %d Paths for %d Ants ===\n", len(paths.AllPaths), ants)
		for i, path := range paths.AllPaths {
			fmt.Printf("\nPath %d (Length: %d):\n", i+1, len(path.Rooms))
			fmt.Printf("Path capacity: %d ants\n", path.CalculateCapacity())
			for j, room := range path.Rooms {
				fmt.Printf("Step %d: Room ID: %d, Flag: %s, Coords: (%d, %d)\n",
					j+1, room.ID, room.Flag, room.X, room.Y)
				if j < len(path.Rooms)-1 {
					fmt.Printf("   ↓   \n")
				}
			}
		}
		fmt.Println("\n=== End of Paths ===")
	*/

	paths, err := Bfs(rooms, ants)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Temporary debug: show paths found
	fmt.Printf("Found %d paths:\n", len(paths.AllPaths))
	for i, path := range paths.AllPaths {
		fmt.Printf("Path %d: ", i)
		for j, room := range path.Rooms {
			if j > 0 {
				fmt.Print("→")
			}
			fmt.Print(room.ID)
		}
		fmt.Printf(" (length: %d)\n", len(path.Rooms)-1)
	}
	fmt.Println()

	// Simulate ant movement
	SimulateAnts(paths, ants, rooms, lines)
}
