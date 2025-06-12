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

	// Find the minimal T that satisfies the constraint
	T := 1
	for {
		totalAnts := 0
		for _, path := range pathsResult.AllPaths {
			antsForPath := T - path.Length + 1
			if antsForPath > 0 {
				totalAnts += antsForPath
			}
		}
		if totalAnts >= numAnts {
			break
		}
		T++
	}

	// Calculate A[i] for each path
	for i, path := range pathsResult.AllPaths {
		antsForPath := T - path.Length + 1
		if antsForPath > 0 {
			distribution[i] = antsForPath
		}
	}

	// Adjust distribution to match exactly numAnts
	totalDistributed := 0
	for _, ants := range distribution {
		totalDistributed += ants
	}
	if totalDistributed > numAnts {
		// Remove excess ants from the longest path
		excess := totalDistributed - numAnts
		distribution[len(distribution)-1] -= excess
	}

	pathsResult.OptimalDistribution = distribution
	return pathsResult, nil
}
