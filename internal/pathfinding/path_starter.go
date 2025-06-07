package pathfinding

import (
	"fmt"
	"lemin/internal/model"
	errhandle "lemin/pkg/errors"
)

func PathStarter(rooms *model.Rooms, numAnts int) (*model.Paths, error) {
	start, err1 := rooms.SearchStart()
	end, err2 := rooms.SearchEnd()
	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("error: no start or end room")
	}

	paths := model.NewPaths()

	// Find all possible paths using DFS
	visited := make(map[int]bool)
	currentPath := []*model.Room{start}
	visited[start.ID] = true

	findAllPaths(start, end, visited, currentPath, paths)

	if len(paths.AllPaths) == 0 {
		return nil, errhandle.FormatError(errhandle.ErrNoPaths)
	}

	// Sort paths by length for optimal distribution
	paths.SortByLength()

	return paths, nil
}

func findAllPaths(curr, end *model.Room, visited map[int]bool, currentPath []*model.Room, paths *model.Paths) {
	if curr == end {
		// Found a path to end
		path := model.Path{
			Rooms:  make([]*model.Room, len(currentPath)),
			Length: len(currentPath) - 1, // -1 because we don't count the start room
		}
		copy(path.Rooms, currentPath)
		paths.AddPath(path)
		return
	}

	for _, next := range curr.Links {
		if !visited[next.ID] {
			visited[next.ID] = true
			currentPath = append(currentPath, next)
			findAllPaths(next, end, visited, currentPath, paths)
			currentPath = currentPath[:len(currentPath)-1] // backtrack
			visited[next.ID] = false
		}
	}
}
