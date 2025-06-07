package main

import (
	"fmt"
	"lemin/model"
)

func SimulateAnts(paths *model.Paths, numAnts int, rooms *model.Rooms, lines []string) {
	// Initialize ants
	ants := make([]*model.Ant, numAnts)
	for i := range ants {
		ants[i] = &model.Ant{
			ID:        i + 1,
			PathIndex: -1, // Will be assigned based on optimal distribution
			StepIndex: 0,  // Start at the start room
			RoomID:    -1, // Will be set when path is assigned
			Active:    false,
			Finished:  false,
		}
	}

	// Distribute ants across paths optimally
	distributeAnts(ants, paths)

	// Track which rooms are occupied (except start and end)
	occupied := make(map[int]bool)

	turn := 0
	finishedAnts := 0
	antStartTurn := make([]int, numAnts) // Track when each ant should start

	// Calculate optimal start times for ants on same path
	for i := 0; i < numAnts; i++ {
		antStartTurn[i] = 1 // Default start on turn 1
		// If multiple ants use same path, delay them
		for j := 0; j < i; j++ {
			if ants[i].PathIndex == ants[j].PathIndex {
				antStartTurn[i] = antStartTurn[j] + 1
			}
		}
	}

	// Print the original input first (as required by the task)
	printOriginalInput(lines)

	// Simulate movement
	for finishedAnts < numAnts {
		turn++
		moves := []string{}

		// Move existing ants first (those already active)
		for _, ant := range ants {
			if ant.Finished || !ant.Active {
				continue
			}

			path := paths.AllPaths[ant.PathIndex]

			// Try to move ant to next room
			if ant.StepIndex+1 < len(path.Rooms) {
				currentRoom := path.Rooms[ant.StepIndex]
				nextRoom := path.Rooms[ant.StepIndex+1]

				// Check if next room is available
				if !occupied[nextRoom.ID] || nextRoom.Flag == "##start" || nextRoom.Flag == "##end" {
					// Free current room (unless it's start or end)
					if currentRoom.Flag != "##start" && currentRoom.Flag != "##end" {
						occupied[currentRoom.ID] = false
					}

					// Move to next room
					ant.StepIndex++
					ant.RoomID = nextRoom.ID
					if nextRoom.Flag != "##start" && nextRoom.Flag != "##end" {
						occupied[nextRoom.ID] = true
					}
					moves = append(moves, fmt.Sprintf("L%d-%d", ant.ID, nextRoom.ID))

					// Check if ant reached the end
					if nextRoom.Flag == "##end" {
						ant.Finished = true
						finishedAnts++
						// Free the room when ant reaches end
						if currentRoom.Flag != "##start" && currentRoom.Flag != "##end" {
							occupied[currentRoom.ID] = false
						}
					}
				}
			}
		}

		// Then try to start new ants (only if it's their turn to start)
		for i, ant := range ants {
			if ant.Finished || ant.Active || turn < antStartTurn[i] {
				continue
			}

			path := paths.AllPaths[ant.PathIndex]

			// Try to move from start to first room
			if ant.StepIndex+1 < len(path.Rooms) {
				nextRoom := path.Rooms[ant.StepIndex+1]
				// Check if next room is occupied (start and end rooms can have multiple ants)
				if !occupied[nextRoom.ID] || nextRoom.Flag == "##start" || nextRoom.Flag == "##end" {
					ant.Active = true
					ant.StepIndex++
					ant.RoomID = nextRoom.ID
					if nextRoom.Flag != "##start" && nextRoom.Flag != "##end" {
						occupied[nextRoom.ID] = true
					}
					moves = append(moves, fmt.Sprintf("L%d-%d", ant.ID, nextRoom.ID))

					// Check if ant reached the end
					if nextRoom.Flag == "##end" {
						ant.Finished = true
						finishedAnts++
					}
				}
			}
		}

		// Print moves for this turn (only if there are moves)
		if len(moves) > 0 {
			for i, move := range moves {
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Print(move)
			}
			fmt.Println()
		}
	}
}

func distributeAnts(ants []*model.Ant, paths *model.Paths) {
	// Better distribution: consider path lengths and distribute optimally
	if len(paths.AllPaths) == 0 {
		return
	}

	// Simple strategy: use shortest paths first, distribute evenly
	// Sort paths by length (already done in BFS)
	if len(paths.AllPaths) >= 2 {
		// Use the two shortest paths
		ants[0].PathIndex = 0 // First ant uses shortest path
		if len(ants) > 1 {
			ants[1].PathIndex = 1 // Second ant uses second shortest path
		}
		if len(ants) > 2 {
			ants[2].PathIndex = 0 // Third ant uses shortest path (delayed)
		}
		// Continue pattern for more ants
		for i := 3; i < len(ants); i++ {
			ants[i].PathIndex = i % len(paths.AllPaths)
		}
	} else {
		// Only one path available
		for _, ant := range ants {
			ant.PathIndex = 0
		}
	}

	// Set initial room for all ants
	for _, ant := range ants {
		ant.RoomID = paths.AllPaths[ant.PathIndex].Rooms[0].ID // Start room
	}
}

func printOriginalInput(lines []string) {
	// Print the original input as required by the task
	for _, line := range lines {
		fmt.Println(line)
	}
	fmt.Println() // Empty line before moves
}
