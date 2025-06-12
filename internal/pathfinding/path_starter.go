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

	// Create flow network
	network := NewFlowNetwork(rooms)

	// Find all possible paths
	pathStrings := network.FindPaths(start.ID, end.ID)
	if len(pathStrings) == 0 {
		return nil, errhandle.FormatError(errhandle.ErrNoPaths)
	}

	// Convert string paths to model.Paths
	pathsResult := model.NewPaths()
	for _, pathStr := range pathStrings {
		path := &model.Path{
			Rooms:  make([]*model.Room, len(pathStr)),
			Length: len(pathStr) - 1,
		}
		for i, roomID := range pathStr {
			room, err := rooms.GetRoomByID(roomID)
			if err != nil {
				return nil, fmt.Errorf("error: invalid room in path: %v", err)
			}
			path.Rooms[i] = room
		}
		pathsResult.AddPath(*path)
	}

	// Sort paths by length
	pathsResult.SortByLength()

	// Calculate optimal distribution
	distribution := make([]int, len(pathStrings))

	// Calculate the optimal number of ants per path
	// We want to maximize parallel movement while minimizing total time
	totalLength := 0
	for _, path := range pathsResult.AllPaths {
		totalLength += path.Length
	}

	// Calculate the base distribution
	// Distribute ants to maximize parallel movement
	remainingAnts := numAnts

	// First, calculate how many ants each path should get based on its length
	// and the total number of paths
	baseAnts := numAnts / len(pathStrings)

	// Distribute base ants
	for i := range distribution {
		distribution[i] = baseAnts
		remainingAnts -= baseAnts
	}

	// Distribute extra ants to maintain staggered movement
	for remainingAnts > 0 {
		// Find the path that will minimize the total time
		bestPath := 0
		bestTime := (distribution[0] + 1) * pathsResult.AllPaths[0].Length
		for i := 1; i < len(distribution); i++ {
			time := (distribution[i] + 1) * pathsResult.AllPaths[i].Length
			if time < bestTime {
				bestTime = time
				bestPath = i
			}
		}
		distribution[bestPath]++
		remainingAnts--
	}

	// Ensure the distribution maintains staggered movement
	// by adjusting the distribution to minimize waiting time
	for i := 1; i < len(distribution); i++ {
		// Calculate the time difference between adjacent paths
		timeDiff := (distribution[i] * pathsResult.AllPaths[i].Length) -
			(distribution[i-1] * pathsResult.AllPaths[i-1].Length)

		// If the time difference is too large, move some ants
		if timeDiff > pathsResult.AllPaths[i].Length {
			move := timeDiff / (2 * pathsResult.AllPaths[i].Length)
			if move > 0 {
				distribution[i] -= move
				distribution[i-1] += move
			}
		}
	}

	// Ensure each path has at least one ant
	for i := range distribution {
		if distribution[i] == 0 {
			// Find the path with the most ants
			maxPath := 0
			for j := 1; j < len(distribution); j++ {
				if distribution[j] > distribution[maxPath] {
					maxPath = j
				}
			}
			// Move one ant from the path with the most ants
			distribution[maxPath]--
			distribution[i]++
		}
	}

	pathsResult.OptimalDistribution = distribution
	return pathsResult, nil
}
