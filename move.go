package main

import (
	"fmt"
	"lemin/model"
)

func SimulateAnts(endID int, path []*model.Room, numAnts int) {
	ants := make([]model.Ant, numAnts)
	for i := range ants {
		ants[i] = model.Ant{ID: i + 1, Pos: -1, Active: false}
	}
	occupied := make([]bool, len(path))
	turn := 0
	for {
		turn++
		moves := []string{}
		done := make(map[int]bool, 0)
		for i := 0; i < numAnts; i++ {
			ant := ants[i]

			if !ant.Active {
				ant.Pos = path[1].ID
				ant.Active = true
				occupied[ant.Pos] = true
				moves = append(moves, fmt.Sprintf("L%d-%d", ant.ID, ant.Pos))
			} else if ant.Pos != endID && !occupied[path[ant.Pos+1].ID] {
				occupied[ant.Pos] = false
				ant.Pos = path[ant.Pos+1].ID
				occupied[ant.Pos] = true
				moves = append(moves, fmt.Sprintf("L%d-%d", ant.ID, ant.Pos))
			} else if ant.Pos == endID {
				done[ant.ID] = true
			}
		}
	}
}
