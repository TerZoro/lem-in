package simulation

import "lemin/internal/model"

func distributeAnts(ants []*model.Ant, paths *model.Paths) {
	// Better distribution: consider path lengths and distribute optimally
	if len(paths.AllPaths) == 0 {
		return
	}

	// Simple strategy: use shortest paths first, then next shortest
	if len(paths.AllPaths) >= 2 {
		// Distribute ants in a pattern: 0,1,0,1,0,1,0,1...
		for i := 0; i < len(ants); i++ {
			if i%2 == 0 {
				ants[i].PathIndex = 0 // Even: use shortest path
			} else {
				ants[i].PathIndex = 1 // Odd: use second shortest path
			}
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
