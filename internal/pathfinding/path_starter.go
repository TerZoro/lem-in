package pathfinding

import (
	"fmt"
	"lemin/internal/model"
	errhandle "lemin/pkg/errors"
)

// PathStarter finds optimal paths and distribution for minimal total turns
func PathStarter(rooms *model.Rooms, numAnts int) (*model.Paths, error) {
	start, err1 := rooms.SearchStart()
	end, err2 := rooms.SearchEnd()
	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("error: no start or end room")
	}

	// Step 1: Find all possible paths
	allPaths := findAllPossiblePaths(start, end)
	if len(allPaths) == 0 {
		return nil, errhandle.FormatError(errhandle.ErrNoPaths)
	}

	// Step 2: Find optimal disjoint path combination
	bestPaths, bestDistribution := findOptimalPathCombination(allPaths, numAnts)

	// Step 3: Create the result with optimal distribution
	result := model.NewPaths()
	for _, path := range bestPaths {
		result.AddPath(*path)
	}

	// Store the optimal distribution for the simulation
	result.OptimalDistribution = bestDistribution

	return result, nil
}

// findAllPossiblePaths uses DFS to find all possible paths from start to end
func findAllPossiblePaths(start, end *model.Room) []*model.Path {
	var allPaths []*model.Path
	visited := make(map[int]bool)
	currentPath := []*model.Room{start}

	var dfs func(curr *model.Room)
	dfs = func(curr *model.Room) {
		if curr == end {
			// Found a complete path
			path := &model.Path{
				Rooms:  make([]*model.Room, len(currentPath)),
				Length: len(currentPath) - 1,
			}
			copy(path.Rooms, currentPath)
			allPaths = append(allPaths, path)
			return
		}

		for _, next := range curr.Links {
			if !visited[next.ID] {
				visited[next.ID] = true
				currentPath = append(currentPath, next)
				dfs(next)
				currentPath = currentPath[:len(currentPath)-1]
				visited[next.ID] = false
			}
		}
	}

	visited[start.ID] = true
	dfs(start)

	return allPaths
}

// findOptimalPathCombination finds the best combination of disjoint paths and ant distribution
func findOptimalPathCombination(allPaths []*model.Path, numAnts int) ([]*model.Path, []int) {
	minTurns := int(^uint(0) >> 1) // Max int
	var bestPaths []*model.Path
	var bestDistribution []int

	// Try single path solutions
	for _, path := range allPaths {
		distribution := []int{numAnts}
		turns := simulateMovement([]*model.Path{path}, distribution)
		if turns < minTurns {
			minTurns = turns
			bestPaths = []*model.Path{path}
			bestDistribution = distribution
		}
	}

	// Try combinations of 2 disjoint paths
	for i, path1 := range allPaths {
		for j, path2 := range allPaths {
			if i >= j {
				continue
			}

			if arePathsDisjoint(path1, path2) {
				// Try different distributions of ants across these 2 paths
				for ants1 := 1; ants1 < numAnts; ants1++ {
					ants2 := numAnts - ants1
					distribution := []int{ants1, ants2}
					turns := simulateMovement([]*model.Path{path1, path2}, distribution)
					if turns < minTurns {
						minTurns = turns
						bestPaths = []*model.Path{path1, path2}
						bestDistribution = distribution
					}
				}
			}
		}
	}

	// Try combinations of 3 disjoint paths (if applicable)
	for i, path1 := range allPaths {
		for j, path2 := range allPaths {
			for k, path3 := range allPaths {
				if i >= j || j >= k {
					continue
				}

				if arePathsDisjoint(path1, path2) && arePathsDisjoint(path1, path3) && arePathsDisjoint(path2, path3) {
					// Try different distributions across 3 paths
					for ants1 := 1; ants1 <= numAnts-2; ants1++ {
						for ants2 := 1; ants2 <= numAnts-ants1-1; ants2++ {
							ants3 := numAnts - ants1 - ants2
							if ants3 > 0 {
								distribution := []int{ants1, ants2, ants3}
								turns := simulateMovement([]*model.Path{path1, path2, path3}, distribution)
								if turns < minTurns {
									minTurns = turns
									bestPaths = []*model.Path{path1, path2, path3}
									bestDistribution = distribution
								}
							}
						}
					}
				}
			}
		}
	}

	return bestPaths, bestDistribution
}

// arePathsDisjoint checks if two paths share any intermediate rooms (excluding start/end)
func arePathsDisjoint(path1, path2 *model.Path) bool {
	// Create a set of intermediate rooms in path1 (excluding start and end)
	path1Rooms := make(map[int]bool)
	for i := 1; i < len(path1.Rooms)-1; i++ {
		path1Rooms[path1.Rooms[i].ID] = true
	}

	// Check if path2 shares any intermediate rooms
	for i := 1; i < len(path2.Rooms)-1; i++ {
		if path1Rooms[path2.Rooms[i].ID] {
			return false
		}
	}

	return true
}

// simulateMovement simulates ant movement to calculate total turns needed
func simulateMovement(paths []*model.Path, distribution []int) int {
	// Initialize ants
	type Ant struct {
		pathIndex int
		stepIndex int
		finished  bool
		startTurn int
	}

	var ants []Ant
	antID := 0

	// Create ants according to distribution
	for pathIdx, numAntsOnPath := range distribution {
		for i := 0; i < numAntsOnPath; i++ {
			ants = append(ants, Ant{
				pathIndex: pathIdx,
				stepIndex: 0,
				finished:  false,
				startTurn: i + 1, // Stagger ants on same path
			})
			antID++
		}
	}

	// Simulate movement
	occupied := make(map[int]bool) // Track occupied intermediate rooms
	turn := 0
	finishedCount := 0

	for finishedCount < len(ants) {
		turn++

		// Try to move each ant
		for i := range ants {
			ant := &ants[i]
			if ant.finished || turn < ant.startTurn {
				continue
			}

			path := paths[ant.pathIndex]

			// Try to advance to next room
			if ant.stepIndex+1 < len(path.Rooms) {
				currentRoom := path.Rooms[ant.stepIndex]
				nextRoom := path.Rooms[ant.stepIndex+1]

				// Check if next room is available
				canMove := false
				if nextRoom.Flag == "##start" || nextRoom.Flag == "##end" {
					canMove = true // Start and end rooms never blocked
				} else {
					canMove = !occupied[nextRoom.ID]
				}

				if canMove {
					// Free current room (if it's an intermediate room)
					if currentRoom.Flag != "##start" && currentRoom.Flag != "##end" {
						occupied[currentRoom.ID] = false
					}

					// Move to next room
					ant.stepIndex++

					// Occupy next room (if it's an intermediate room)
					if nextRoom.Flag != "##start" && nextRoom.Flag != "##end" {
						occupied[nextRoom.ID] = true
					}

					// Check if reached end
					if nextRoom.Flag == "##end" {
						ant.finished = true
						finishedCount++
					}
				}
			}
		}
	}

	return turn
}
