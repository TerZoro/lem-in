package simulation

import "lemin/internal/model"

func distributeAnts(ants []*model.Ant, paths *model.Paths) {
	if len(paths.AllPaths) == 0 {
		return
	}

	// Use the optimal distribution calculated by the pathfinding algorithm
	if len(paths.OptimalDistribution) > 0 && len(paths.AllPaths) >= 2 {
		// Use simple alternating pattern: 0,1,0,1,0,1...
		// This ensures proper ordering: path0 gets ants 1,3,5... path1 gets ants 2,4,6...
		for i := 0; i < len(ants); i++ {
			pathIndex := i % len(paths.AllPaths)
			ants[i].PathIndex = pathIndex
			ants[i].RoomID = paths.AllPaths[pathIndex].Rooms[0].ID // Start room
		}
	} else if len(paths.OptimalDistribution) > 0 {
		// Single path case - all ants use the same path
		for i := 0; i < len(ants); i++ {
			ants[i].PathIndex = 0
			ants[i].RoomID = paths.AllPaths[0].Rooms[0].ID // Start room
		}
	} else {
		// Fallback to simple strategy if no optimal distribution available
		if len(paths.AllPaths) >= 2 {
			for i := 0; i < len(ants); i++ {
				if i%2 == 0 {
					ants[i].PathIndex = 0
				} else {
					ants[i].PathIndex = 1
				}
			}
		} else {
			for _, ant := range ants {
				ant.PathIndex = 0
			}
		}

		// Set initial room for fallback case
		for _, ant := range ants {
			ant.RoomID = paths.AllPaths[ant.PathIndex].Rooms[0].ID
		}
	}
}
