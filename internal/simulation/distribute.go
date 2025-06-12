package simulation

import "lemin/internal/model"

func distributeAnts(ants []*model.Ant, paths *model.Paths) {
	if len(paths.AllPaths) == 0 {
		return
	}

	// Use the optimal distribution calculated by the flow network algorithm
	if len(paths.OptimalDistribution) > 0 {
		antIndex := 0
		// Distribute ants according to the optimal distribution
		for pathIndex, numAntsOnPath := range paths.OptimalDistribution {
			for i := 0; i < numAntsOnPath; i++ {
				if antIndex < len(ants) {
					ants[antIndex].PathIndex = pathIndex
					ants[antIndex].RoomID = paths.AllPaths[pathIndex].Rooms[0].ID // Start room
					antIndex++
				}
			}
		}
	} else {
		// Fallback to simple strategy if no optimal distribution available
		for i := 0; i < len(ants); i++ {
			pathIndex := i % len(paths.AllPaths)
			ants[i].PathIndex = pathIndex
			ants[i].RoomID = paths.AllPaths[pathIndex].Rooms[0].ID // Start room
		}
	}
}
